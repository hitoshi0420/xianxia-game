# xianxia-game - Xianxia Cultivation Browser Game

A text-based xianxia (Chinese cultivation fantasy) game played in the browser. Start as a mortal and cultivate your way to immortality.

## Features
- Character creation with stats (spiritual root, cultivation level)
- Turn-based cultivation system
- Alchemy, artifact crafting, and skill training
- Combat system against demons and rival cultivators
- Save/load game progress
- Browser-based UI (no installation needed)

## Tech Stack
- Go (backend server)
- HTML/CSS/JavaScript (browser frontend)
- JSON file persistence

## Project Structure
`
xianxia-game/
  main.go              # Entry point (starts server on :8080)
  internal/
    server/            # HTTP server & game engine
  go.mod               # Go module definition
  xianxia_save.json    # Save file
`

## Quick Start
`ash
go run main.go
# Or run the pre-built binary:
./xianxia-game.exe
`

Open http://localhost:8080 in your browser.

## Gameplay
- Choose your spiritual root (affects cultivation speed)
- Practice meditation to gain cultivation points
- Learn techniques, craft pills and artifacts
- Break through cultivation bottlenecks
- Face tribulations to advance to higher realms

## License
MIT
