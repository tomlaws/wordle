package protocol

import (
	"encoding/json"
	"fmt"
)

func NewProtocol(
	registry map[MessageType]func() Payload,
) *Protocol {
	return &Protocol{
		registry: registry,
	}
}

func (p *Protocol) wrapMessage(payload Payload) (*Message, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("wrap failed: %w", err)
	}
	msg := Message{
		Type:    payload.MessageType(),
		Payload: data,
	}
	return &msg, nil
}

func (p *Protocol) unwrapMessage(msg *Message) (Payload, error) {
	constructor, ok := p.registry[msg.Type]
	if !ok {
		return nil, fmt.Errorf("unknown message type: %s", msg.Type)
	}

	payload := constructor()
	if err := json.Unmarshal(msg.Payload, payload); err != nil {
		return nil, fmt.Errorf("unmarshal failed for type %s: %w", msg.Type, err)
	}
	return payload, nil
}

func (p *Protocol) WrapChannel(ch chan json.RawMessage) chan Payload {
	wrapped := make(chan Payload)
	go func() {
		defer close(wrapped)
		for payload := range wrapped {
			msg, err := p.wrapMessage(payload)
			if err != nil {
				// Handle error (e.g., log it)
				continue
			}
			data, err := json.Marshal(msg)
			if err != nil {
				// Handle error (e.g., log it)
				continue
			}
			ch <- data
		}
	}()
	return wrapped
}

func (p *Protocol) UnwrapChannel(ch chan json.RawMessage) chan Payload {
	unwrapped := make(chan Payload)
	go func() {
		defer close(unwrapped)
		for rawMsg := range ch {
			var msg Message
			if err := json.Unmarshal([]byte(rawMsg), &msg); err != nil {
				// Handle error (e.g., log it)
				continue
			}
			payload, err := p.unwrapMessage(&msg)
			if err != nil {
				// Handle error (e.g., log it)
				continue
			}
			unwrapped <- payload
		}
	}()
	return unwrapped
}
