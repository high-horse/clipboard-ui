package clip

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ClipboardItem struct {
	ID        int
	Content   string
	Timestamp time.Time
}

type ClipboardUIEvent struct {
	ClipAction Action
	Payload    string
}

type ClipboardManager struct {
	mu         sync.Mutex
	items      []ClipboardItem
	currentId  int
	nextItemID int
	ctx        context.Context
}

func NewClipboardManager() *ClipboardManager {
	return &ClipboardManager{}
}

func (cm *ClipboardManager) SetContext(ctx context.Context) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.ctx = ctx
}

func (cm *ClipboardManager) GetContext() context.Context {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.ctx
}

func (cm *ClipboardManager) Add(content string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	item := ClipboardItem{
		ID:        cm.nextItemID,
		Content:   content,
		Timestamp: time.Now(),
	}

	cm.nextItemID++
	cm.items = append([]ClipboardItem{item}, cm.items...) // make newest first
	cm.currentId = item.ID
	if cm.ctx != nil {
		runtime.EventsEmit(cm.ctx, "new-content", map[string]string{
			"content": content,
		})
		log.Println("event emitted 'new-content'")
	}
}

func (cm *ClipboardManager) Select(id int) bool {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for i, item := range cm.items {
		if item.ID == id {
			// cm.currentId = id
			// return true
			cm.items = append([]ClipboardItem{item}, append(cm.items[:i], cm.items[i+1:]...)...)
			cm.currentId = id
			return true
		}
	}
	return false
}

func (cm *ClipboardManager) GetCurrentContext() string {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for _, item := range cm.items {
		if cm.currentId == item.ID {
			return item.Content
		}
	}

	return ""
}

func (cm *ClipboardManager) GetHistory() []ClipboardItem {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	return cm.items
}

func (cm *ClipboardManager) ClearHistory() {
	cm.mu.Lock()
	cm.mu.Unlock()

	cm.items = []ClipboardItem{}
	cm.currentId = -1
	cm.nextItemID = 0
}
