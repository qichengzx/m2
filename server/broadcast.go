package server

import (
	"github.com/hashicorp/memberlist"
)

type memberlistBroadcast struct {
	msg    []byte
	notify chan<- struct{}
}

func (b *memberlistBroadcast) Invalidates(other memberlist.Broadcast) bool {
	return false
}

func (b *memberlistBroadcast) Message() []byte {
	return b.msg
}

func (b *memberlistBroadcast) Finished() {
	if b.notify != nil {
		close(b.notify)
	}
}
