# Ink Mark

A local-first, single-user markdown writing app built with **Go + Wails** (desktop backend), **SvelteKit + TypeScript** (frontend), and **Tailwind CSS** (styling).

---

## Architecture

```md
ink-mark/
├── .devcontainer/ # VS Code dev container config
│ ├── devcontainer.json
│ ├── docker-compose.yml
│ ├── Dockerfile
│ └── post-create.sh
├── frontend/ # SvelteKit + TypeScript + Tailwind app
│ ├── src/
│ │ ├── lib/
│ │ │ └── wailsjs/ # Auto-generated Wails JS bindings (Go → JS bridge)
│ │ └── routes/ # SvelteKit pages
│ └── build/ # Compiled frontend (embedded into Go binary)
├── app.go # Go app struct — bound methods exposed to frontend
├── main.go # Wails entry point
└── wails.json # Wails project config
```

---

## Getting started

1. Open this repo in VS Code
2. When prompted, click **"Reopen in Container"**
3. The `post-create.sh` script will automatically install Go modules and frontend dependencies

### Running in dev mode

```bash
cd /workspace/ink-mark
wails dev
```

### Dev servers

| Service                | URL                    | Description                             |
| ---------------------- | ---------------------- | --------------------------------------- |
| **Wails dev server**   | http://localhost:34115 | Full app with Go ↔ JS bridge (use this) |
| **Vite frontend only** | http://localhost:5173  | UI only, no Go backend calls            |

> **Note:** Always use `http://localhost:34115` during development — it has the full Wails bridge so Go-bound methods (like `Greet()`) work from the browser.

### Building for production

```bash
cd /workspace/ink-mark
wails build
```

The compiled desktop binary will be output to `build/bin/`.
