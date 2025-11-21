# Go Site Clone

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)
![Wails](https://img.shields.io/badge/Wails-3.0--alpha-DF5B00?style=flat)
![Nuxt](https://img.shields.io/badge/Nuxt-4.1.2-00DC82?style=flat&logo=nuxt.js)
![License](https://img.shields.io/badge/License-MIT-green?style=flat)

A powerful website cloning and local deployment tool built with Wails3 + Nuxt.js

[Features](#-features) • [Quick Start](#-quick-start) • [Usage](#-usage) • [Architecture](#-architecture) • [Configuration](#-configuration)

</div>

---

## 📖 Overview

Go Site Clone is a powerful website cloning tool that can completely capture all resources of a website (HTML, CSS, JavaScript, images, videos, etc.) and support local deployment. It uses a real browser engine for resource capture, equipped with comprehensive anti-detection mechanisms and browser fingerprint spoofing system.

### ✨ Highlights

- 🚀 **Real Browser Engine** - Based on Chrome DevTools Protocol, simulating real user behavior
- 🎭 **Fingerprint Spoofing** - Canvas fingerprint, timezone, language, UserAgent comprehensive disguise
- 🔄 **Smart Resource Management** - Auto deduplication, same-domain filtering, categorized storage
- 📦 **Complete Site Cloning** - Support for static resources and dynamic content
- 🖥️ **Cross-Platform** - One-click packaging for Windows, macOS, Linux
- 💫 **Modern UI** - Nuxt.js + Ant Design Vue responsive interface
- 🌐 **Local Site Service** - Integrated Nginx for one-click local mirror site startup

---

## 🎯 Features

### Resource Capture

- ✅ Complete HTML page cloning
- ✅ Batch CSS stylesheet download
- ✅ JavaScript file capture
- ✅ Image resources (JPG, PNG, GIF, SVG, WebP)
- ✅ Video files (MP4, WebM, etc.)
- ✅ Fonts and other static resources
- ✅ Smart link traversal (auto-track internal links)

### Anti-Detection

- 🎭 Canvas fingerprint randomization
- 🌍 Timezone simulation (global timezone support)
- 🗣️ Language and locale spoofing
- 🔐 Custom UserAgent
- 📍 Geolocation coordinate spoofing
- 🕵️ Stealth plugin integration (bypass anti-bot detection)

### File Management

- 📁 Auto-categorized storage by domain
- 🔍 Resource type classification
- 📊 Real-time download progress
- 🗂️ Local site list view
- 📂 One-click folder opening

### Local Deployment

- 🌐 Integrated Nginx server
- ⚡ Quick start local mirror sites
- 🔧 Custom port configuration
- 📝 Site configuration management

---

## 🚀 Quick Start

### Requirements

- Go 1.24 or higher
- Node.js 18+ and npm/pnpm
- Wails CLI v3

### Installation

#### 1. Install Wails CLI

```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
```

#### 2. Clone Repository

```bash
git clone https://github.com/zhuy1228/go-site-clone.git
cd go-site-clone
```

#### 3. Install Go Dependencies

```bash
go mod download
```

#### 4. Install Frontend Dependencies

```bash
cd frontend
npm install
# or use pnpm
pnpm install
```

### Development Mode

```bash
# Start dev server with hot-reload
wails3 dev
```

The app will start in development mode with auto-reload for both frontend and backend changes.

### Production Build

```bash
# Build production version
wails3 build

# Output in build/ directory
```

---

## 📚 Usage

### Basic Workflow

1. **Start Application**
   ```bash
   wails3 dev  # or run the built executable
   ```

2. **Enter Target Website URL**
   - Input the website URL you want to clone
   - Example: `https://example.com`

3. **Get Resource List**
   - Click "Get Resources" button
   - System automatically analyzes and lists all resources

4. **Download Site**
   - Select resource types to download
   - Click "Download" to start cloning
   - View real-time download progress

5. **Local Preview**
   - Select a site from download list
   - Click "Open Folder" to view files
   - Or use integrated Nginx service for local deployment

### Configuration

Edit `config.yaml` to customize settings:

```yaml
# Application name
appName: "Go Site Clone"

# Local service port
port: 6997

# WebSocket connection URL
wsUrl: "106.12.33.188:6996"

# STUN server address
stunUrl: "stun:106.12.33.188:3478"

# API service URL
apiUrl: "http://106.12.33.188:6996"

# Website file save directory
siteFileDir: "www"
```

### API Methods

The application exposes the following methods for frontend:

#### `GetResources(url string)`
Get all resource lists for specified URL

**Parameters:**
- `url`: Target website URL

**Returns:** `ResourcesList` object with categorized resources

#### `DownloadSite(uri string, resources ResourcesList)`
Download website resources to local

**Parameters:**
- `uri`: Website URL
- `resources`: Resource list object

**Returns:** `bool` - Success status

#### `GetDownloadList()`
Get list of downloaded websites

**Returns:** `[]FileDir` - File directory list

#### `OpenSiteFileDir(path string)`
Open folder for specified website

**Parameters:**
- `path`: Relative path

**Returns:** `bool` - Operation success status

---

## 🏗️ Architecture

### Tech Stack

**Backend (Go)**
- [Wails v3](https://wails.io/) - Application framework
- [go-rod](https://github.com/go-rod/rod) - Browser automation
- [go-rod/stealth](https://github.com/go-rod/stealth) - Anti-detection
- [goquery](https://github.com/PuerkitoBio/goquery) - HTML parsing
- [go-sqlite3](https://github.com/mattn/go-sqlite3) - Database
- [yaml.v3](https://github.com/go-yaml/yaml) - Configuration

**Frontend (Nuxt.js)**
- [Nuxt 4](https://nuxt.com/) - Vue framework
- [Vue 3](https://vuejs.org/) - Reactive framework
- [Ant Design Vue](https://antdv.com/) - UI component library
- [Chart.js](https://www.chartjs.org/) - Charts
- [Wails Runtime](https://wails.io/) - Frontend-backend communication

### Project Structure

```
go-site-clone/
├── app.go                      # Application main logic
├── main.go                     # Program entry
├── config.yaml                 # Configuration file
├── go.mod                      # Go module dependencies
│
├── config/                     # Configuration loader
│   └── index.go
│
├── services/                   # Business services
│   ├── site_service.go        # Website cloning core
│   └── nginx_service.go       # Nginx service management
│
├── libs/                       # Core libraries
│   ├── browser.go             # Browser pool management
│   ├── chrome.go              # Chrome instance creation
│   ├── sqlite.go              # Database operations
│   ├── pack.go                # Packaging utilities
│   └── browser-fingerprint/   # Browser fingerprint spoofing
│       ├── index.go
│       ├── canvas.go          # Canvas fingerprint
│       └── timezone.go        # Timezone simulation
│
├── types/                      # Type definitions
│   ├── fingerprint_type.go    # Fingerprint parameter types
│   ├── task_type.go           # Task types
│   └── nginx_template.go      # Nginx templates
│
├── utils/                      # Utility functions
│   ├── file.go                # File download/management
│   └── index.go               # Environment checks
│
├── frontend/                   # Nuxt.js frontend
│   ├── nuxt.config.ts         # Nuxt configuration
│   ├── package.json           # Frontend dependencies
│   ├── app/                   # Application code
│   │   ├── app.vue           # Root component
│   │   ├── components/       # Vue components
│   │   ├── pages/            # Page routes
│   │   ├── layouts/          # Layout templates
│   │   └── plugins/          # Plugins
│   └── bindings/              # Go binding types
│
├── user-data/                  # Browser user data
├── www/                        # Website file storage
├── plugin/                     # Plugins
│   └── nginx/                 # Nginx plugin
└── build/                      # Build configuration
    ├── darwin/                # macOS packaging config
    ├── linux/                 # Linux packaging config
    └── windows/               # Windows packaging config
```

---

## ⚙️ Configuration

### Browser Fingerprint Config

```go
type BrowserFingerprintParams struct {
    Canvas      bool   // Randomize Canvas fingerprint
    TimeZone    string // Timezone, e.g., "America/New_York"
    Language    string // Language, e.g., "en-US"
    UserAgent   string // Custom UserAgent
    GeoLocation string // Geolocation, e.g., "40.7128,-74.0060"
}
```

---

## 🔧 Development

### Adding Features

1. **Backend Services**
   - Create new service in `services/`
   - Implement business logic
   - Expose methods in `app.go` for frontend

2. **Frontend UI**
   - Add new pages in `frontend/app/pages/`
   - Add components in `frontend/app/components/`
   - Use Wails Runtime to call backend methods

### Debugging

```bash
# View detailed logs
wails3 dev -v

# Build with verbose output
wails3 build -v

# Clean build cache
wails3 build -clean
```

---

## 📦 Build & Deploy

### Windows

```bash
wails3 build -platform windows/amd64
```

Executable in `build/bin/`

### macOS

```bash
wails3 build -platform darwin/universal
```

Generates `.app` bundle

### Linux

```bash
wails3 build -platform linux/amd64
```

Supports `.deb`, `.rpm`, AppImage formats

---

## 🤝 Contributing

Issues and Pull Requests are welcome!

1. Fork the repository
2. Create feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Submit Pull Request

---

## 📄 License

This project is licensed under the MIT License - see [LICENSE](LICENSE) file

---

## 🙏 Acknowledgments

- [Wails](https://wails.io/) - Excellent Go + Web framework
- [go-rod](https://go-rod.github.io/) - Powerful browser automation
- [Nuxt.js](https://nuxt.com/) - Elegant Vue framework
- [Ant Design Vue](https://antdv.com/) - Beautiful UI components

---

## 📞 Contact

- Author: zhuy1228
- Project: [https://github.com/zhuy1228/go-site-clone](https://github.com/zhuy1228/go-site-clone)

---

<div align="center">

**If this project helps you, please give it a ⭐️**

Made with ❤️ by zhuy1228

**[中文文档](README_CN.md)**

</div>
