package server

import (
	"github.com/hashicorp/memberlist"
)

type memberlistBroadcast struct {
	msg    []byte
	node   string
	notify chan<- struct{}
}

func (b *memberlistBroadcast) Invalidates(other memberlist.Broadcast) bool {
	mb, ok := other.(*memberlistBroadcast)
	if !ok {
		return false
	}

	return b.node == mb.node
}

func (b *memberlistBroadcast) Name() string {
	return b.node
}

func (b *memberlistBroadcast) Message() []byte {
	return b.msg
}

func (b *memberlistBroadcast) Finished() {
	select {
	case b.notify <- struct{}{}:
	default:
	}
}
