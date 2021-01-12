package syncmap

import (
	"errors"
	"sync"
)

type Stdmap struct {
	db sync.Map
}

func New() *Stdmap {
	return &Stdmap{
		db: sync.Map{},
	}
}

func (s *Stdmap) Set(key, value []byte) error {
	s.db.Store(string(key), value)
	return nil
}

func (s *Stdmap) Get(key []byte) ([]byte, error) {
	if val, ok := s.db.Load(string(key)); ok {
		return val.([]byte), nil
	}
	return nil, errors.New("empty")
}

func (s *Stdmap) Delete(key []byte) error {
	s.db.Delete(string(key))
	return nil
}

func (s *Stdmap) Close() error {
	return nil
}
