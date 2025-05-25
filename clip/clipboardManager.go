package clip

import (
	"context"
	"encoding/binary"
	"errors"
	"log"
	"sort"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.etcd.io/bbolt"
)


type ClipboardManager struct {
	db *bbolt.DB
	maxLen int
	ctx context.Context
}


type CopiedContent struct {
	Key   uint64 `json:"key"`
	Value string `json:"value"`
}



const BUCKET = "clipboard_bucket"

func NewClipboardManager(path string) (*ClipboardManager, error){
	db, err := bbolt.Open(path, 0666, &bbolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		return nil, err
	}

	return &ClipboardManager{
		db: db,
		// maxLen: 200,
		maxLen: 5,
	}, nil
}


func (cm *ClipboardManager) SetContext(ctx context.Context) {
	cm.ctx = ctx
}

func (cm *ClipboardManager) GetContext() context.Context {
	return cm.ctx
}


func (cm *ClipboardManager) Add(content string) error {

	return cm.db.Update(func(tx *bbolt.Tx) error{
		b, err := tx.CreateBucketIfNotExists([]byte(BUCKET))
		if err != nil {
			return err
		}

		id,_ := b.NextSequence()
		key := itob(id)

		if err := b.Put(key, []byte(content)); err != nil {
			return err
		}
		runtime.EventsEmit(cm.ctx, "new-content", map[string]CopiedContent{
			"content": {
				Key: id,
				Value: content,
			},
		})

		// trim over max len
		if b.Stats().KeyN > cm.maxLen {
			var keys [][]byte

			_ = b.ForEach(func(k, _ []byte) error {
				keys = append(keys, k)
				return nil
			})

			sort.Slice(keys, func(i int, j int) bool{
				return binary.BigEndian.Uint64(keys[i]) < binary.BigEndian.Uint64(keys[j])
			})

			for i := 0; i < len(keys)-cm.maxLen; i++ {
				runtime.EventsEmit(cm.ctx, "remove-content", map[string]int64 {
					"index": btoi(keys[i]) ,
				})
				_ = b.Delete(keys[i])
			}

		}
		return nil
	})
}



func (cm *ClipboardManager) GetAll() ([]CopiedContent, error) {
	var results []string
	var cpiedContents []CopiedContent

	err := cm.db.View(func(tx *bbolt.Tx) error{
		b := tx.Bucket([]byte(BUCKET))
		if b == nil {
			return errors.New("Clipboard bucket not found")
		}

		c := b.Cursor()

		for k, v:=  c.Last(); k != nil; k,v = c.Prev() {
			results = append(results, string(v))
			cpiedContents = append(cpiedContents, CopiedContent{
				Key:   binary.BigEndian.Uint64(k),
				Value: string(v),
			})
		}
		return nil
	})

	log.Println("returning ", cpiedContents)
	return cpiedContents, err
}

func itob(v uint64) []byte {
    b := make([]byte, 8)
    binary.BigEndian.PutUint64(b, v)
    return b
}

func btoi(v []byte) int64 {
	if(len(v) < 8){
		padded := make([]byte, 8)
		copy(padded[8-len(v):], v)
		v = padded
	}
	return int64(binary.BigEndian.Uint64(v))
}
