package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Enabled reports whether agent-loop logging is active.
// Controlled by the INKMARK_DEV_LOGS=1 environment variable.
func Enabled() bool {
	return os.Getenv("INKMARK_DEV_LOGS") == "1"
}

// OpenLogFile opens (or creates and appends to) logs/<name>.log in the
// project root — resolved relative to the working directory at startup.
// Returns nil if INKMARK_DEV_LOGS is not "1", so callers can always
// check for nil and skip writing.
func OpenLogFile(name string) *os.File {
	if !Enabled() {
		return nil
	}

	dir := filepath.Join("logs")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "devlog: failed to create logs dir: %v\n", err)
		return nil
	}

	path := filepath.Join(dir, name+".log")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "devlog: failed to open %s: %v\n", path, err)
		return nil
	}

	return f
}

// WriteLine writes a single JSON-Lines entry to w.
// Safe to call with a nil writer — becomes a no-op.
func WriteLine(w io.Writer, level, source, msg string) {
	if w == nil {
		return
	}
	ts := time.Now().UTC().Format(time.RFC3339Nano)
	// Minimal hand-rolled JSON — avoids pulling encoding/json for a one-liner.
	fmt.Fprintf(w, "{\"ts\":%q,\"level\":%q,\"source\":%q,\"msg\":%q}\n", ts, level, source, msg)
}
