package server

import (
	"encoding/json"
	"net/http"
)

func (srv *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	val := r.Form.Get("val")
	if key == "" {
		http.Error(w, "error key is empty", http.StatusOK)
		return
	}

	byteKey := []byte(key)
	byteVal := []byte(val)
	err := srv.delegate.Set(byteKey, byteVal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(Payload{
		Action: "set",
		Data: struct {
			Key   []byte
			Value []byte
		}{Key: byteKey, Value: byteVal},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	srv.broadcastChan <- &memberlistBroadcast{
		msg:    b,
		notify: nil,
	}
}

func (srv *Server) DelHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	if key == "" {
		http.Error(w, "error key is empty", http.StatusOK)
		return
	}
	byteKey := []byte(key)
	err := srv.delegate.Delete(byteKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(Payload{
		Action: "del",
		Data: struct {
			Key   []byte
			Value []byte
		}{Key: byteKey, Value: nil}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	broadcastQueue.QueueBroadcast(&memberlistBroadcast{
		msg:    b,
		notify: nil,
	})
}

func (srv *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	if key == "" {
		http.Error(w, "error key is empty", http.StatusOK)
		return
	}
	val, err := srv.delegate.Get([]byte(key))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(val)
}
