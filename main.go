package main

import (
	"context"
	"embed"
	"log"
	"os"
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
	log.SetFlags(log.Ldate| log.Ltime | log.Lshortfile)
	// Create an instance of the app structure
	app := NewApp()
	clipboard := clip.NewClipboardManager()
	// clipboard.SetContext(app.ctx)
	
	go watchClipboard(clipboard)
	select{
		case <-exitSignal :
			log.Println("Exiting due to clipboard init faliure ")
			return 
		default:
	}
	

	// Create application with options
	err := wails.Run(&options.App{
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
		Bind: []interface{}{
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
		log.Println("error occured in ",err)
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