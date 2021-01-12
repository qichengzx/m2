package server

import (
	"github.com/hashicorp/memberlist"
	"os"
)

var (
	hostname       string
	broadcastQueue *memberlist.TransmitLimitedQueue
)

func init() {
	hostname, _ = os.Hostname()
}
