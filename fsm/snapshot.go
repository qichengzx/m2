package fsm

import (
	"github.com/hashicorp/raft"
)

type snapshot struct{}

func (s snapshot) Persist(_ raft.SnapshotSink) error {
	return nil
}

func (s snapshot) Release() {}

func newSnapshot() (raft.FSMSnapshot, error) {
	return &snapshot{}, nil
}
