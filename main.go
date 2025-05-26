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
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/getlantern/systray"

	// hook "github.com/moutend/go-hook"
	// hook "github.com/robotn/gohook"
	cb "golang.design/x/clipboard"
	"golang.design/x/hotkey"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed frontend/src/assets/images/logo-universal.png
var trayIcon []byte

var (
	clipboard *clip.ClipboardManager
	app        *App
	wailsCtx   context.Context
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(trayIcon)
	systray.SetTitle("G-Clipboard")
	systray.SetTooltip("G-Clipboard running.")

	mShow := systray.AddMenuItem("Show", "Show the main menu")
	mHide := systray.AddMenuItem("Hide", "Hide the menu item")
	systray.AddSeparator()
	mClear := systray.AddMenuItem("Clear clipboard History", "Clear all clipboard history.")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the app")

	// App setup
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("Failed to get user home directory:", err)
	}
	app := NewApp()
	configPath := filepath.Join(homeDir, ".config", "clipboard-ui")
	if err = os.MkdirAll(configPath, 0775); err != nil {
		log.Fatalln("Failed to create config directory:", err)
	}
	dbPath := filepath.Join(configPath, "clipboard_bucket.db")
	log.Println("db path is ", dbPath)

	clipboard, err = clip.NewClipboardManager(dbPath)
	if err != nil {
		log.Fatalln("Failed to initialize clipboard manager:", err)
	}

	// Start Wails app in a goroutine
	go func() {
		err = wails.Run(&options.App{
			Title:  "ui-clipboard",
			Width:  1024,
			Height: 768,
			AssetServer: &assetserver.Options{
				Assets: assets,
			},
			BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
			OnStartup: func(ctx context.Context) {
				wailsCtx = ctx
				app.startup(ctx)
				clipboard.SetContext(ctx)
				go watchClipboard(ctx, clipboard)
				go listenForHotkey(ctx)
			},
			OnBeforeClose: app.beforeClose,
			Bind: []any{
				app,
				clipboard,
			},
		})
		if err != nil {
			log.Println("Error:", err.Error())
		}
	}()

	// Tray menu handler
	go func() {
		for {
			select {
			case <-mShow.ClickedCh:
				if wailsCtx != nil {
					runtime.Show(wailsCtx)
				}
			case <-mHide.ClickedCh:
				if wailsCtx != nil {
					runtime.Hide(wailsCtx)
				}
			case <-mClear.ClickedCh:
				if clipboard != nil {
					// err := clipboard.Clear()
					// if err != nil {
					// 	log.Println("Failed to clear clipboard:", err)
					// } else {
						log.Println("Clipboard history cleared")
					// }
				}
			case <-mQuit.ClickedCh:
				systray.Quit()
				os.Exit(0)
			}
		}
	}()

}

func onExit() {

}

func main_() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("Failed to get user home directory:", err)
	}
	// go systray.Run(onReady, onExit)

	// Create an instance of the app structure
	app := NewApp()
	configPath := filepath.Join(homeDir, ".config", "clipboard-ui")
	if err = os.MkdirAll(configPath, 0775); err != nil {
		log.Fatalln("Failed to create config directory:", err)
	}
	dbPath := filepath.Join(configPath, "clipboard_bucket.db")
	log.Println("db path is ", dbPath)

	clipboard, err := clip.NewClipboardManager(dbPath)
	if err != nil {
		log.Fatalln("Failed to initialize clipboard manager:", err)
	}

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "ui-clipboard",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		// StartHidden:      true,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			clipboard.SetContext(ctx)

			go watchClipboard(ctx, clipboard)
			go listenForHotkey(ctx)
		},
		OnBeforeClose: app.beforeClose,
		Bind: []any{
			app,
			clipboard,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func watchClipboard(ctx context.Context, cBoard *clip.ClipboardManager) {
	err := cb.Init()
	if err != nil {
		log.Println("error occured in ", err)
		return
	}

	log.Println("Clipboard watcher initialized.")
	board := cb.Watch(context.Background(), cb.FmtText)

	for {
		select {
		case <-ctx.Done():
			log.Println("App is shutting down, stopping clipboard watcher.")
			return

		case data := <-board:
			log.Println("new content detected ", string(data))
			cBoard.Add(string(data))
		}
	}
}

func listenForHotkey(ctx context.Context) {
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyV)
	err := hk.Register()
	if err != nil {
		log.Printf("Failed to register hotkey: %v", err)
		return
	}
	log.Println("Global hotkey Ctrl+Alt+V registered")
	for {
		select {
		case <-ctx.Done():
			return
		case <-hk.Keydown():
			log.Println("Hotkey pressed: Ctrl+Alt+V")
			runtime.Show(ctx)
		}
	}
}
