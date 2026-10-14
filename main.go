package main

import (
	"embed"
	"io"
	"log"
	"os"

	"ink-mark/internal/logging"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

// Version is injected at build time via -ldflags "-X main.Version=x.y.z".
// Falls back to "dev" when built without the flag.
var Version = "dev"

//go:embed all:frontend/build
var assets embed.FS

func main() {
	// When INKMARK_DEV_LOGS=1, tee all Go log output to logs/go.log so
	// agents can read it without watching a terminal.
	if f := logging.OpenLogFile("go"); f != nil {
		defer f.Close()
		log.SetOutput(io.MultiWriter(os.Stderr, f))
		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	}

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "ink-mark",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 247, G: 246, B: 243, A: 0},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
