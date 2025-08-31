# Multiplayer Wordle Game Design

## Rules
- Two players compete to guess the same hidden word.
- The game consists of 12 rounds in total (6 turns per player).
- Players take turns alternately.
- Each player has up to 60 seconds to submit a guess per turn. If time runs out, the turn is skipped.
- After each guess, feedback is provided and visible to both players.
- The first player to correctly guess the word wins immediately.
- If neither player guesses the word within 12 rounds, the game ends in a tie.

### Rule Justification
The rules are designed for simplicity and fast-paced gameplay. Limiting each player to 6 turns and 60 seconds per guess keeps matches short and engaging, reducing downtime and making the game accessible to new players. Alternating turns and immediate feedback ensure fairness and maintain excitement.

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
WebSocket is selected for this game because it supports bidirectional communication and is compatible with all modern browsers, making it ideal for real-time multiplayer interactions.

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
            "nickname": "Tom",
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
            },
        ]
    }
}
```

This structured format ensures clear, extensible communication for all game events.