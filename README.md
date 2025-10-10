# Multiplayer Wordle

A real-time, multiplayer Wordle game built with Go, featuring queue-based matchmaking and a WebSocket protocol.

<img src="https://i.imgur.com/EwnRwVv.jpeg" width="480"/>

## Features
- Multiplayer: Two players compete to guess the same hidden word.
- Real-time feedback: Both players see each other's guesses and feedback.
- Queue-based matchmaking: Players are matched automatically for quick games.
- Turn-based: Players alternate turns, each with a time limit.
- WebSocket protocol: Efficient, bidirectional communication between client and server.
- Configurable word list and game settings.

## Getting Started

### Prerequisites
- Go 1.20 or later

### Installation
```sh
git clone https://github.com/tomlaws/wordle.git
cd wordle
go mod download
```

### Running the Server
```sh
go run cmd/server/main.go
```
or to provide a custom configuration
```sh
go run -ldflags="-X main.Port=8080 -X main.MaxGuesses=6 -X main.WordListPath=assets/words.txt -X main.ThinkTime=60" cmd/server/main.go
```
#### Configuration
- **Port:** The server listens on port 8080 by default.
- **Max Guesses:** The maximum number of guesses in a game is 6 by default.
- **Word List:** The default word list is located at `assets/words.txt`.
- **Think Time:** Each player has 60 seconds per turn by default.

### Running the Console Client
```sh
go run cmd/client/main.go
```

### Running the Web Client

#### Setup
Install dependencies:
```sh
cd web
npm install
```

#### Start in Development Mode
Launch the development server with hot reloading:
```sh
npm run dev
```
The app will be available at [http://localhost:5173](http://localhost:5173) by default.

#### Build for Production
Generate an optimized production build:
```sh
npm run build
```
The output will be in the `web/build` directory.

## Usage
- Start the server and client as above.
- The client will connect to the server, join the matchmaking queue, and start a game when matched.
- Enter your guesses when prompted. Each guess must be a valid 5-letter word.
- The first player to guess the word wins. If neither guesses correctly in 6 rounds, the game is a tie.

## Design
See `GAME_DESIGN.md` for full details on architecture, game flow, and design decisions.

## Standalone Version
### Running the Standalone
```sh
go run cmd/standalone/main.go
```
or to provide a custom configuration
```sh
go run -ldflags="-X main.MaxGuesses=6 -X main.WordListPath=assets/words.txt" cmd/standalone/main.go
```
#### Configuration
- **Max Guesses:** The maximum number of guesses is 6 by default.
- **Word List:** The default word list is located at `assets/words.txt`.

## Acknowledgments
- Inspired by [Wordle](https://www.nytimes.com/games/wordle/index.html).
- Built with Go and the Gorilla WebSocket library.
