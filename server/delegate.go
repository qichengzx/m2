package server

import (
	"encoding/json"
	"github.com/qichengzx/m2/storage"
	"log"
	"sync"
)

type delegate struct {
	sync.RWMutex
	meta        []byte
	state       []byte
	remoteState []byte
	storage.DB
}

func (d *delegate) NodeMeta(limit int) []byte {
	d.Lock()
	defer d.Unlock()

	return d.meta
}

func (d *delegate) NotifyMsg(msg []byte) {
	if len(msg) == 0 {
		return
	}

	var payload Payload
	if err := json.Unmarshal(msg, &payload); err != nil {
		log.Println("error:", err)
		return
	}

	switch payload.Action {
	case "set":
		d.Set(payload.Data.Key, payload.Data.Value)

	case "del":
		d.Delete(payload.Data.Key)
	}
}

func (b *delegate) GetBroadcasts(overhead, limit int) [][]byte {
	return broadcastQueue.GetBroadcasts(overhead, limit)
}

func (d *delegate) LocalState(join bool) []byte {
	d.RLock()
	defer d.RUnlock()

	return d.state
}

func (d *delegate) MergeRemoteState(s []byte, join bool) {
	d.Lock()
	defer d.Unlock()

	d.remoteState = s
}
