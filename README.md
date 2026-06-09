# Can You Run It? (Go)

This project is a repository for learning Golang. It is a desktop application built with [Wails](https://wails.io/) that helps users determine if their computer meets the system requirements to run specific PC games.

## Features

- **System Information Gathering**: Detects the user's hardware specifications (CPU, GPU, RAM, OS) using native system calls.
- **Game Requirements Checking**: Fetches minimum and recommended system requirements for games (integrating with Steam/SteamSpy data).
- **Performance Dashboard**: Provides a visual representation of how your hardware stacks up against the game's requirements.
- **Modern UI**: A sleek, reactive frontend to display game requirements and system hardware comparisons.

## Tech Stack

- **Backend**: Go (Golang)
- **Frontend**: Web technologies (HTML, CSS, JavaScript/TypeScript) via Wails
- **Data Sources**: SteamSpy / Steam API for game data

## Getting Started

### Prerequisites

- [Go](https://go.dev/doc/install)
- [Node.js](https://nodejs.org/en/download/)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)

### Installation & Running

1. Clone the repository and navigate into the project directory:
   ```bash
   git clone <repository_url>
   cd canurunitgo
   ```

2. Run the application in development mode (this will automatically start the frontend server and the Go backend):
   ```bash
   wails dev
   ```

3. Build the application for production:
   ```bash
   wails build
   ```

## Learning Goals

As a learning project, the main objectives include:
- Exploring Go's system-level operations for hardware detection.
- Building cross-platform desktop applications using Go and web technologies.
- Structuring a Go project effectively (e.g., using `/internal`, `/pkg` directories).
- Fetching and parsing data from external APIs (like SteamSpy).
