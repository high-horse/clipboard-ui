# Clipboard UI

A cross-platform clipboard manager desktop application built with [Wails](https://wails.io/) (Go backend) and [Vue 3 + TypeScript](https://vuejs.org/) frontend.

## Features

- Modern UI with Vue 3 and Vite
- Clipboard history and management
- Native system integration via Go
- Hot reload for rapid development
- Cross-platform: Linux, Windows, macOS

## Prerequisites

- [Go](https://golang.org/) (see `go.mod` for version)
- [Node.js](https://nodejs.org/) (LTS recommended)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation/)

## Getting Started

1. **Clone the repository:**
   ```bash
   git clone <your-repo-url>
   cd ui-clipboard
   ```

2. **Install frontend dependencies:**
   ```bash
    cd frontend
    npm install
    cd ..
   ```

3. **Check Wails setup:**
   ```bash
   wails doctor
   ```
   
## Development
To run in live development mode with hot reload:
   ```bash
   wails dev
   ```

## Building for Production
   ```bash
   wails build
   ```
   The output binaries will be in the build/bin or target directory.


## Development
   ```bash
   .
├── [app.go](http://_vscodecontentref_/1)                # Main Go application logic
├── clip/                 # Clipboard management Go code
├── frontend/             # Vue 3 + TypeScript frontend
│   ├── src/              # Vue components, assets, types
│   └── wailsjs/          # Wails-generated JS bindings
├── build/                # Build assets and output
├── [wails.json](http://_vscodecontentref_/2)            # Wails project configuration
└── [README.md](http://_vscodecontentref_/3)             # This file
   ```

## Configuration

Edit `wails.json` to configure project settings. See the [Wails documentation](https://wails.io/docs/reference/project-config) for more details.