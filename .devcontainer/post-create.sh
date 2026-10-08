#!/usr/bin/env bash
set -euo pipefail

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  Ink Mark — post-create setup"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# ─── Go backend ───────────────────────────────────────────────────────────
if [ -f "/workspace/ink-mark/backend/go.mod" ]; then
  echo "▶ Downloading Go modules..."
  cd /workspace/ink-mark/backend && go mod download
fi

# ─── Frontend ─────────────────────────────────────────────────────────────
if [ -f "/workspace/ink-mark/frontend/package.json" ]; then
  echo "▶ Installing frontend dependencies..."
  cd /workspace/ink-mark/frontend && pnpm install
fi

echo ""
echo "✅ Dev environment is ready!"
echo ""
echo "  Frontend   → http://localhost:3000  (pnpm dev)"
echo "  Backend    → http://localhost:8080  (go run ./...)"
echo "  PostgreSQL → localhost:5432  (inkmark / inkmark_dev)"
echo ""

