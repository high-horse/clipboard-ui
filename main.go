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

	cb "golang.design/x/clipboard"
	"golang.design/x/hotkey"
	"fyne.io/systray"
	"sync/atomic"
	"time"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed frontend/src/assets/images/logo-universal.png
var trayIcon []byte

var (
	clipboard *clip.ClipboardManager
	app       *App
	wailsCtx  context.Context
)
var lastCopiedValue atomic.Value // for thread-safe access

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// App setup
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("Failed to get user home directory:", err)
	}
	
	app = NewApp()
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

	// Run Wails app on main thread (this is essential for proper window rendering)
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
			
			// Start background processes
			go watchClipboard(ctx, clipboard)
			go listenForHotkey(ctx)
			go setupSystemTrayAlternative(ctx)
		},
		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			// Don't actually close the app, just hide it
			runtime.Hide(ctx)
			return true // Prevent closing
		},
		Bind: []any{
			app,
			clipboard,
		},
	})
	
	if err != nil {
		log.Println("Error:", err.Error())
	}
}

func setupSystemTrayAlternative(ctx context.Context) {
    go func() {
        var mLastClip *systray.MenuItem

        systray.Run(func() {
            systray.SetIcon(trayIcon)
            systray.SetTitle("G-Clipboard")
            systray.SetTooltip("G-Clipboard running")

            mShow := systray.AddMenuItem("Show", "Show window")
            mLastClip = systray.AddMenuItem("Last: (empty)", "Show last clipboard item")
            mClear := systray.AddMenuItem("Clear Clipboard", "Clear clipboard history")
            mQuit := systray.AddMenuItem("Quit", "Quit app")

            // Update the tray menu/tooltip periodically
            go func() {
                for {
                    v := lastCopiedValue.Load()
                    text := "(empty)"
                    if v != nil && v.(string) != "" {
                        text = v.(string)
                    }
                    mLastClip.SetTitle("Last: " + text)
                    systray.SetTooltip("Last copied: " + text)
                    time.Sleep(500 * time.Millisecond)
                }
            }()

            for {
                select {
                case <-mShow.ClickedCh:
                    runtime.Show(ctx)
                case <-mLastClip.ClickedCh:
                    v := lastCopiedValue.Load()
                    if v != nil {
                        log.Println("Last clipboard item:", v.(string))
                    }
                case <-mClear.ClickedCh:
                    if clipboard != nil {
                        // clipboard.Clear() if implemented
                        log.Println("Clipboard cleared.")
                    }
                case <-mQuit.ClickedCh:
                    runtime.Quit(ctx)
                    return
                }
            }
        }, func() {})
    }()
}

// Alternative approach: Create system tray using a different method
// You'll need to replace this with a system tray library that works better with Wails
// For now, this is a placeholder that demonstrates the concept
func setupSystemTrayAlternative_(ctx context.Context) {
	// Option 1: Use fyne system tray (install: go get fyne.io/systray)
	// Option 2: Use a cross-platform system tray library like energye/systray
	// Option 3: Create a simple HTTP server for tray-like functionality
	
	// For demonstration, I'll show you how to implement basic tray functionality
	// without the problematic getlantern/systray library
	
	log.Println("System tray alternative setup - implement your preferred tray library here")
	
	// You can implement a simple menu system or use runtime.Menu* functions
	// to create application menus instead of system tray
	setupAppMenu(ctx)
}

func setupAppMenu(ctx context.Context) {
	// Create application menu instead of system tray
	// This will show in the application's menu bar
	
	// Example of creating application menus (cross-platform)
	go func() {
		// You can create a simple background service that listens for
		// specific events or implement this using Wails' menu system
		
		// For now, just log that the alternative tray is running
		log.Println("Alternative tray system running - window can be controlled via hotkey")
	}()
}

func watchClipboard(ctx context.Context, cBoard *clip.ClipboardManager) {
	err := cb.Init()
	if err != nil {
		log.Println("error occurred in clipboard init:", err)
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
			log.Println("new content detected:", string(data))
			cBoard.Add(string(data))
			lastCopiedValue.Store(string(data))
			// systray.SetTooltip("Last: " + string(data)) 
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
	
	defer hk.Unregister()
	log.Println("Global hotkey Ctrl+Shift+V registered")
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-hk.Keydown():
			log.Println("Hotkey pressed: Ctrl+Shift+V")
			// Toggle window visibility
			toggleWindow(ctx)
		}
	}
}

func toggleWindow(ctx context.Context) {
	// Simple toggle - you might want to track window state more precisely
	runtime.Show(ctx)
	runtime.WindowUnminimise(ctx)
	runtime.WindowSetAlwaysOnTop(ctx, true)
	
	// Remove always on top after a brief moment
	go func() {
		// time.Sleep(100 * time.Millisecond)
		runtime.WindowSetAlwaysOnTop(ctx, false)
	}()
}