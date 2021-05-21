package server

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/hashicorp/raft"
)

type Server struct {
	raft *raft.Raft
	db   *badger.DB
}

func New(raft *raft.Raft, db *badger.DB) *Server {
	return &Server{
		raft: raft,
		db:   db,
	}
}
