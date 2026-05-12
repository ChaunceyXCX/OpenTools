# ZTools (Wails Port)

<div align="center">

<img src="./.github/assets/icon.png" alt="ZTools Logo" width="120">

**A High-Performance, Extensible Application Launcher and Plugin Platform**

_Open Source Implementation of uTools | Supports macOS, Windows, and Linux_

[![License](https://img.shields.io/github/license/lzx8589561/ZTools)](./LICENSE)
[![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Windows%20%7C%20Linux-blue)](https://github.com/ChaunceyXCX/OpenTools)

English | [简体中文](./README.md)

</div>

---

## ✨ Features

- 🚀 **Quick Launch** - Pinyin search, regex matching, history tracking, pinned apps
- 🧩 **Plugin System** - JS script plugins and Go Service plugins with full API
- 📋 **Clipboard Management** - History tracking, search, auto-monitoring
- 🎨 **Theme Customization** - CSS variable-driven dark theme
- ⚡ **High Performance** - bbolt database, native Go performance, minimal memory
- 🌍 **Cross-Platform** - Native support for macOS, Windows, and Linux
- 🔒 **Data Isolation** - Independent plugin data storage
- 🛠️ **Developer Friendly** - Hot reload frontend, type-safe Go Bindings
- ⚙️ **Tech Stack** - Wails v3 + Go 1.25 + Vue 3 + TypeScript
- 🧩 **Browser Automation** - ZBrowser (chromedp)
- 🔌 **MCP Server** - AI tool integration (scan/launch/clipboard)
- ☁ **WebDAV Sync** - Multi-device data sync

## 🚀 Quick Start

### Installation

> This is the **Wails port** of ZTools (`wails` branch), replacing Electron with Go + WebView.
> The original Electron version is on the `main` branch.

#### Linux Dependencies

```bash
# Ubuntu/Debian
sudo apt install libgtk-3-dev libwebkit2gtk-4.1-dev libx11-dev

# Fedora
sudo dnf install gtk3-devel webkit2gtk4.1-devel libX11-devel

# Arch
sudo pacman -S gtk3 webkit2gtk-4.1 libx11
```

#### Build from Source

```bash
# Clone (wails branch)
git clone -b wails https://github.com/ChaunceyXCX/OpenTools.git
cd OpenTools

# Install frontend deps
cd frontend && pnpm install && cd ..

# Development
wails3 dev

# Build
wails3 build                    # Current platform
wails3 build -platform linux    # Linux
wails3 build -platform darwin   # macOS
wails3 build -platform windows  # Windows
```

### Usage

1. Launch the app, press `Alt+Z` to toggle the main window
2. Type app names or commands to search
3. `↑` `↓` to navigate, `Enter` to launch, `Esc` to hide

### Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| `Alt+Z` | Toggle window |
| `Esc` | Hide window |
| `↑` `↓` | Navigate results |
| `Enter` | Launch selected app |

## 🧩 Plugin System

The Wails port provides full API compatibility with upstream ZTools plugins.

### Plugin Types

- **Go Service Plugins** - Compile-time registered, high performance
- **JS Script Plugins** - goja runtime, supports `window.exports[code].args.enter()`
- **ZPX Package Plugins** - Full `.zpx` format compatibility (asar+gzip)

### Plugin API

Plugins access capabilities via `window.ztools.*`:

| API | Description |
|-----|-------------|
| `ztools.db.*` | Database CRUD (sync + Promise) |
| `ztools.clipboard.*` | Clipboard read/write |
| `ztools.shell.*` | Shell command execution |
| `ztools.showNotification()` | System notification |
| `ztools.showToast()` | Toast notification |
| `ztools.setSubInput()` | Sub-input for plugin search |
| `ztools.hideWindow()` | Hide window |
| `ztools.simulateKeyboardTap()` | Keyboard simulation |

## 🗂 Project Structure

```
OpenTools/
├── main.go                  # Go entry point
├── internal/                # Go backend
│   ├── api/                 # Wails Binding services
│   ├── core/                # Core modules
│   │   ├── database/        # bbolt database
│   │   ├── clipboard/       # Clipboard monitor
│   │   ├── scanner/         # App scanner
│   │   ├── launcher/        # App launcher
│   │   ├── mcp/             # MCP Server
│   │   ├── httpserver/      # HTTP API
│   │   ├── zbrowser/        # Browser automation
│   │   ├── sync/            # WebDAV sync
│   │   └── translation/     # Offline translation
│   ├── plugin/              # Plugin runtime
│   └── native/              # Native capabilities
├── frontend/                # Vue 3 frontend
│   ├── src/
│   │   ├── components/      # UI components
│   │   ├── stores/          # Pinia state management
│   │   └── style.css        # Global styles
│   └── bindings/            # Auto-generated Wails bindings
└── build/                   # Build configs
```

## 📦 Tech Stack

| Component | Technology |
|-----------|-----------|
| Desktop Framework | Wails v3 (alpha.90) |
| Backend | Go 1.25 |
| Frontend | Vue 3 + TypeScript |
| State Management | Pinia |
| Search | Fuse.js |
| Database | bbolt (LMDB compatible) |
| JS Engine | goja (plugin runtime) |
| Browser Automation | chromedp |
| Clipboard | atotto/clipboard |
| WebDAV | gowebdav |
| Packaging | Wails native build |

## 🔄 Branch Strategy

| Branch | Description |
|--------|-------------|
| `main` | Upstream sync branch, matches original ZTools |
| `wails` | Wails port branch (current) |

The `wails` branch periodically merges from `main` to stay up-to-date with upstream changes.

## 📄 License

[MIT](./LICENSE)
