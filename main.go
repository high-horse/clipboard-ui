package main

import (
	"context"
	"embed"
	"log"
	"os"
	"path/filepath"
	"ui-clipboard/clip"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	cb "golang.design/x/clipboard"
)

//go:embed all:frontend/dist
var assets embed.FS

var cpBoard = make(chan os.Signal, 1)
var exitSignal = make(chan struct{})

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("Failed to get user home directory:", err)
	}

	// Create an instance of the app structure
	app := NewApp()
	configPath := filepath.Join(homeDir, ".config", "clipboard-ui")
	if err = os.MkdirAll(configPath, 0775); err != nil {
		log.Fatalln("Failed to create config directory:", err)
	}
	dbPath := filepath.Join(configPath, "clipboard_bucket.db")

	clipboard, err := clip.NewClipboardManager(dbPath)
	if err != nil {
		log.Fatalln("Failed to initialize clipboard manager:", err)
	}

	go watchClipboard(clipboard)
	select {
	case <-exitSignal:
		log.Println("Exiting due to clipboard init faliure ")
		return
	default:
	}

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "ui-clipboard",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			clipboard.SetContext(ctx)
		},
		Bind: []any{
			app,
			clipboard,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func watchClipboard(cBoard *clip.ClipboardManager) {
	err := cb.Init()
	if err != nil {
		log.Println("error occured in ", err)
		exitSignal <- struct{}{}
		return
	}
	log.Println("Clipboard watcher initialized.")
	ch := cb.Watch(context.Background(), cb.FmtText)
	for data := range ch {
		log.Println("new content detected ", string(data))
		cBoard.Add(string(data))
	}
}
