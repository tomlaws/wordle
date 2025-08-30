# Multiplayer Wordle Game Design

## Rules
- Two players compete to guess the same hidden word.
- The game consists of 12 rounds in total (6 turns per player).
- Players take turns alternately.
- Each player has up to 60 seconds to submit a guess per turn. If time runs out, the turn is skipped.
- After each guess, feedback is provided and visible to both players.
- The first player to correctly guess the word wins immediately.
- If neither player guesses the word within 12 rounds, the game ends in a tie.

## Form of Matchmaking
- Room-based Matchmaking
- Queue-based Matchmaking

## Message Protocol
- Websocket (Bidirectional communication channel)