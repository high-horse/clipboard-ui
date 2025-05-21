dev:
	wails dev -tags webkit2_41

install-bolt-browser:	
	go install github.com/br0xen/boltbrowser@latest

view-clipboard: 
	boltbrowser /home/camle/.config/clipboard-ui/clipboard_bucket.db
