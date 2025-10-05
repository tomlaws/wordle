# Multiplayer Wordle Game Design

## Table of Contents
1. [Rules](#rules)
    - [Rule Justification](#rule-justification)
2. [Matchmaking Mechanism Comparison](#matchmaking-mechanism-comparison)
3. [Protocol](#protocol)
    - [Communication Protocol Comparison](#communication-protocol-comparison)
4. [Concurrency](#concurrency)
5. [Message Format](#message-format)
    - [Example Messages](#example-messages)
6. [Player Authentication](#player-authentication)
7. [Error Handling](#error-handling)

---

## Rules
- Two players compete to guess the same hidden word.
- The game consists of 6 rounds in total (3 turns per player).
- Players take turns alternately.
- Each player has up to 60 seconds to submit a guess per turn. If time runs out, the turn is skipped.
- After each guess, feedback is provided and visible to both players.
- The first player to correctly guess the word wins immediately.
- If neither player guesses the word within 6 rounds, the game ends in a tie.

### Rule Justification
The rules are designed for simplicity and fast-paced gameplay. Limiting each player to 3 turns and 60 seconds per guess keeps matches short and engaging, reducing downtime and making the game accessible to new players. Alternating turns and immediate feedback ensure fairness and maintain excitement.

---

## Matchmaking Mechanism Comparison

| Feature                | Room-based Matchmaking         | Queue-based Matchmaking        |
|------------------------|-------------------------------|-------------------------------|
| Player Control         | Players can create/join rooms | Players are matched automatically |
| Waiting Time           | Depends on room availability  | Usually minimized by queue    |
| Customization          | High (room settings possible) | Low (standardized matching)   |
| Social Interaction     | Can invite friends            | Random opponents              |
| Implementation Complexity | Moderate (room management logic) | Low (simple queue logic)         |

**Design Choice:**  
Queue-based matchmaking is favored for its simplicity and faster implementation, allowing for efficient player matching without the overhead of room management.

---

## Protocol

### Communication Protocol Comparison

| Feature                      | HTTP                         | WebSocket                    | gRPC                         |
|------------------------------|------------------------------|------------------------------|------------------------------|
| Communication Model          | Request/Response             | Full-duplex, Bidirectional   | Bidirectional Streaming      |
| Latency                      | Higher (per request)         | Low (persistent connection)  | Low (efficient binary proto) |
| Real-time Support            | Limited (polling/long-poll)  | Native                       | Native                       |
| Message Format               | Text (JSON, HTML, etc.)      | Text/Binary                  | Binary (Protocol Buffers)    |
| Browser Compatibility        | Universal (all browsers)     | Universal (all modern browsers) | Limited (requires HTTP/2, not natively supported in browsers) |
| Ease of Implementation       | Easiest                      | Moderate                     | Most complex                 |

**Design Choice:**  
WebSocket is selected for this game because it supports bidirectional communication and is compatible with all modern browsers, making it ideal for real-time multiplayer interactions and provides extensibility for future enhancements, such as adding a web frontend.

---

## Concurrency

The server uses Go's concurrency model by using two goroutines per connected client: one for reading messages from the WebSocket and another for writing messages to it. This allows the server to handle incoming and outgoing communication independently, ensuring that a slow or blocked client does not stall the entire game loop.

The client also uses separate goroutines for handling incoming and outgoing WebSocket messages. This ensures console/UI update and user input do not block each other.

---

## Message Format

Messages between client and server use JSON over WebSocket. Each message includes a `type` field to indicate its purpose and a `payload` object for relevant data.

### Example Messages

#### Player Guess

```json
{
    "type": "guess",
    "payload": {
        "word": "apple"
    }
}
```

#### Feedback

```json
{
    "type": "feedback",
    "payload": {
        "player": {
            "id": "3d6a2e36-30c0-4812-89a6-39bcb1b6edc2",
            "nickname": "Tom"
        },
        "round": 1,
        "feedback": [
            {
                "letter": "a",
                "position": 0,
                "match_type": 1 // 1 indicates the letter is in the word but in the wrong position
                },
            {
                "letter": "p",
                "position": 1,
                "match_type": 0 // 0 indicates letter is not in the word
                },
            {
                "letter": "p",
                "position": 2,
                "match_type": 0
            },
            {
                "letter": "l",
                "position": 3,
                "match_type": 0
            },
            {
                "letter": "e",
                "position": 4,
                "match_type": 2 // 2 indicates correct letter in correct position
            }
        ]
    }
}
```

This structured format ensures clear, extensible communication for all game events.

---

## Player Authentication

Players are required to enter a username only when connecting to the game. The server only establishes a connection if a username is provided, ensuring that every player has an identifiable display name. To keep the implementation simple and memory-efficient, the system does not check for duplicate usernames. Instead, each player is assigned a unique UUID upon connection, which allows the client to distinguish between the local player and their opponent, regardless of username duplication.

While this approach simplifies the design and implementation, it also means that player identity can be easily forged by providing any username. Since there is no authentication or duplicate username check, the system should not be used for scenarios where secure or trusted player identification is required.

---

## Error Handling

- **Validation Errors:**  
    When a client sends an invalid message (e.g., malformed JSON, missing required fields, or invalid guess), the server responds with an explicit error message. This message includes a `type` so the client can display or handle the error gracefully.

    ```json
    {
            "type": "invalid_word",
            "payload": {
                "player": {
                    "id": "3d6a2e36-30c0-4812-89a6-39bcb1b6edc2",
                    "nickname": "Tom",
                },
                "round": 1,
                "word": "aeiou"
            }
    }
    ```

- **Critical Errors:**  
    For unrecoverable or critical errors, the server closes the WebSocket connection. This approach keeps the implementation simple and avoids complex error recovery logic on the client side.

