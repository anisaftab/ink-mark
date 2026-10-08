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

# ─── Frontend ─────────────────────────────────────────────────────────────
if [ -f "/workspace/ink-mark/frontend/package.json" ]; then
  echo "▶ Installing frontend dependencies..."
  cd /workspace/ink-mark/frontend && npm install
fi

echo ""
echo "✅ Dev environment is ready!"
echo ""
echo "  Run:  cd /workspace/ink-mark && wails dev"
echo ""
echo "  Wails dev server  → http://localhost:34115  (full Go bridge)"
echo "  Vite frontend     → http://localhost:5173   (UI only)"
echo ""


