export interface Payload {
    MessageType(): string;
}

export interface Message {
    type: string;
    payload?: Payload;
}

export class Protocol {
    private payloadRegistry: Map<string, () => Payload>;
    constructor(
        registerPayloads: Map<string, () => Payload> = new Map<string, () => Payload>()
    ) {
        this.payloadRegistry = registerPayloads;
    }

    createMessage(payload: Payload): Message {
        return {
            type: payload.MessageType(),
            payload
        };
    }

    parseMessage(data: Message | Object): Payload {
        console.log('Parsing message:', data);
        if (typeof data !== 'object' || data === null || !('type' in data)) {
            throw new Error(`Invalid message format`);
        }
        const PayloadClass = this.payloadRegistry.get(data.type);
        if (!PayloadClass) {
            throw new Error(`Unknown message type: ${data.type}`);
        }
        return Object.assign(PayloadClass(), data.payload);
    }
}
