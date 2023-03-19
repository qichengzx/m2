package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/qichengzx/m2/fsm"
)

func (s *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	val := r.Form.Get("val")
	if key == "" || val == "" {
		http.Error(w, "error key or val is empty", http.StatusOK)
		return
	}

	payload := fsm.Payload{
		OP:    fsm.CMDSET,
		Key:   key,
		Value: val,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	applyFuture := s.raft.Apply(data, 500*time.Millisecond)
	if err := applyFuture.Error(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, ok := applyFuture.Response().(*fsm.ApplyResponse)
	if !ok {
		w.Write([]byte("error raft response"))
		return
	}

	w.Write([]byte("ok"))
}

func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	if key == "" {
		http.Error(w, "error key is empty", http.StatusOK)
		return
	}
	var value = make([]byte, 0)
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		value, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(value)
}

func (s *Server) DelHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	if key == "" {
		http.Error(w, "error key is empty", http.StatusOK)
		return
	}
	payload := fsm.Payload{
		OP:  fsm.CMDDEL,
		Key: key,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	applyFuture := s.raft.Apply(data, 500*time.Millisecond)
	if err := applyFuture.Error(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, ok := applyFuture.Response().(*fsm.ApplyResponse)
	if !ok {
		w.Write([]byte("error raft response"))
		return
	}
	w.Write([]byte("ok"))
}
