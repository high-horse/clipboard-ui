dev:
	wails dev

dev-tags:
	wails dev -tags webkit2_41

install-bolt-browser:	
	go install github.com/br0xen/boltbrowser@latest

view-clipboard: 
	boltbrowser /home/camle/.config/clipboard-ui/clipboard_bucket.db

build:
	wails build
	
build-clean:
	wails build -clean


build-tags:
	wails build -tags webkit2_41
