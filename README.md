# Tic-Tac-Toe

Local LAN multiplayer tic-tac-toe. One person hosts the server; anyone on the same network can join by opening a URL in their browser. No internet connection or cloud hosting required.

## How to play

**Prerequisites:** Go 1.26+ and Node.js 24 (LTS)

This project uses [nvm](https://github.com/nvm-sh/nvm) to pin the Node version. Switch to the correct version before building:

```bash
nvm install   # installs the version from .nvmrc (Node 24) if not already present
nvm use       # activates it for the current shell
```

Then build and run:

```bash
make build   # builds the frontend then compiles the Go binary
./tic-tac-claude
```

Or to build and run in one step:

```bash
make run
```

The terminal will print two URLs:

```
  Tic-Tac-Toe
  ──────────────────────────────────────────
  Player 1 (host)  http://localhost:8080
  Player 2 (LAN)   http://192.168.x.x:8080
  ──────────────────────────────────────────
  Both players can use the LAN IP — even the host.
  Ctrl+C to stop the server.
```

- **Player 1** opens `http://localhost:8080` (or the LAN IP) on the host machine.
- **Player 2** opens the LAN IP URL on any device on the same network — phone, laptop, tablet.

The game starts automatically when both players are connected. Player 1 is always X, Player 2 is always O. After a win or draw, either player can click **Play Again** to reset, or **Quit** to disconnect and leave the server running for a new game.

## Project structure

```
├── main.go       # Entry point — LAN IP detection, server lifecycle
├── server.go     # HTTP routes, embeds the built frontend
├── hub.go        # WebSocket hub — player slots, game state, broadcast
├── handler.go    # Per-connection WebSocket read loop
├── game.go       # Pure game logic — board, moves, win/draw detection
├── messages.go   # JSON message structs for the WebSocket protocol
└── frontend/
    └── src/
        ├── types.ts           # TypeScript types mirroring the Go message protocol
        ├── useGame.ts         # WebSocket hook — connection and all game state
        ├── App.tsx            # Root component
        └── components/
            ├── Board.tsx      # 3×3 grid with click handling and win highlights
            ├── StatusBar.tsx  # Turn indicator and status messages
            └── Controls.tsx   # Post-game Restart / Quit buttons
```

## Tech stack

- **Backend:** Go, [gorilla/websocket](https://github.com/gorilla/websocket)
- **Frontend:** React, TypeScript, Tailwind CSS, Vite
- The compiled frontend is embedded directly into the Go binary at build time — no separate file server needed.

---

Built with [Claude](https://claude.ai)
