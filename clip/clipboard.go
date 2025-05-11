package clip

import (
	"sync"
	"time"
)


type ClipboardItem struct{
	ID int
	Content string
	Timestamp time.Time
}

type ClipboardUIEvent struct {
	ClipAction Action
	Payload string
}

type ClipboardManager struct {
	mu sync.Mutex
	items []ClipboardItem
	currentId int
	nextItemID int
}

func NewClipboardManager() *ClipboardManager {
	return &ClipboardManager{}
}

func (c *ClipboardManager) Add (content string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	item := ClipboardItem{
		ID: c.nextItemID,
		Content: content,
		Timestamp: time.Now(),
	}
	
	c.nextItemID++
	c.items = append([]ClipboardItem{item}, c.items...) // make newest first
	c.currentId = item.ID
}

func (c *ClipboardManager) Select(id int) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	for i, item := range c.items {
		if item.ID == id {
			// c.currentId = id
			// return true
			c.items = append([]ClipboardItem{item}, append(c.items[:i], c.items[i+1:]...)...)
			c.currentId = id
			return true
		}
	}	
	return false
}

func (c *ClipboardManager) GetCurrentContext() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	for _, item := range c.items {
		if c.currentId == item.ID {
			return item.Content
		}
	}
	
	return ""
}

func (c *ClipboardManager) GetHistory() []ClipboardItem{
	c.mu.Lock()
	defer c.mu.Unlock()
	
	return c.items
}

func (c *ClipboardManager) ClearHistory() {
	c.mu.Lock()
	c.mu.Unlock()
	
	c.items = []ClipboardItem{}
	c.currentId = -1
	c.nextItemID = 0
}