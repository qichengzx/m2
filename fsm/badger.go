package fsm

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/hashicorp/raft"
	"io"
	"os"
)

const (
	CMDSET = "SET"
	CMDDEL = "DEL"
)

type badgerFSM struct {
	db *badger.DB
}

func (b badgerFSM) set(key string, value interface{}) error {
	var data = make([]byte, 0)
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if data == nil || len(data) <= 0 {
		return nil
	}

	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), data)
	})
}

func (b badgerFSM) get(key string) (interface{}, error) {
	var data []byte
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		data, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}

		return nil
	})

	return data, err
}

func (b badgerFSM) delete(key string) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

type Payload struct {
	OP    string
	Key   string
	Value interface{}
}

type ApplyResponse struct {
	Error error
	Data  interface{}
}

func (b badgerFSM) Apply(log *raft.Log) interface{} {
	switch log.Type {
	case raft.LogCommand:
		var payload = Payload{}
		if err := json.Unmarshal(log.Data, &payload); err != nil {
			fmt.Fprintf(os.Stderr, "error marshalling payload %s\n", err.Error())
			return nil
		}

		switch payload.OP {
		case CMDSET:
			return &ApplyResponse{
				Error: b.set(payload.Key, payload.Value),
				Data:  payload.Value,
			}
		case CMDDEL:
			return &ApplyResponse{
				Error: b.delete(payload.Key),
				Data:  nil,
			}
		}
	}
	fmt.Fprintf(os.Stderr, "raft log command type:%s\n", raft.LogCommand)
	return nil
}

func (b badgerFSM) Snapshot() (raft.FSMSnapshot, error) {
	return newSnapshot()
}

func (b badgerFSM) Restore(rc io.ReadCloser) error {
	return nil
}

func NewBadger(badgerDB *badger.DB) raft.FSM {
	return &badgerFSM{
		db: badgerDB,
	}
}
