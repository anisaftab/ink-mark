#!/usr/bin/env bash
set -euo pipefail

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  Ink Mark — post-create setup"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# ─── Go modules ───────────────────────────────────────────────────────────
if [ -f "/workspace/ink-mark/go.mod" ]; then
  echo "▶ Downloading Go modules..."
  cd /workspace/ink-mark && go mod download
fi

# ─── Wails CLI ────────────────────────────────────────────────────────────
echo "▶ Installing Wails CLI..."
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# ─── Goose CLI ────────────────────────────────────────────────────────────
echo "▶ Installing Goose CLI..."
go install github.com/pressly/goose/v3/cmd/goose@latest

# ─── Frontend ─────────────────────────────────────────────────────────────
if [ -f "/workspace/ink-mark/frontend/package.json" ]; then
  echo "▶ Installing frontend dependencies..."
  cd /workspace/ink-mark/frontend && npm install
fi

echo ""
echo "✅ Dev environment is ready!"
echo ""
echo "  UI only (hot-reload):  cd /workspace/ink-mark/frontend && npm run dev"
echo "  Full dev (Go bridge):  cd /workspace/ink-mark && wails dev"
echo ""
echo "  Vite frontend     → http://localhost:5173   (UI only, no Go bridge)"
echo "  Wails dev server  → http://localhost:34115  (full Go bridge)"
echo ""

