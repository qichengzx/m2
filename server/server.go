package server

import (
	"github.com/hashicorp/memberlist"
	"github.com/qichengzx/m2/storage"
	"github.com/qichengzx/m2/storage/rocksdb"
	"github.com/qichengzx/m2/storage/syncmap"
	"log"
	"strconv"
	"sync"
)

type Server struct {
	delegate *delegate
	name     string
}

func New(db string) *Server {
	var storage storage.DB
	switch db {
	case "rocksdb":
		storage = rocksdb.New("")
		break

	case "syncmap":
		storage = syncmap.New()
		break
	default:
		panic(db + " not exist")
	}

	return &Server{
		name: "",
		delegate: &delegate{
			RWMutex:     sync.RWMutex{},
			meta:        []byte{},
			state:       []byte{},
			remoteState: []byte{},
			DB:          storage,
		},
	}
}

type Payload struct {
	Action string
	Data   struct {
		Key   []byte
		Value []byte
	}
}

func (srv *Server) Start(port int, members []string) error {
	c := memberlist.DefaultLocalConfig()
	c.BindPort = port + 10000
	c.Delegate = srv.delegate
	c.Events = &event{}
	c.Name = hostname + "-" + strconv.Itoa(port)

	m, err := memberlist.Create(c)
	if err != nil {
		return err
	}
	if len(members) > 0 {
		_, err := m.Join(members)
		if err != nil {
			return err
		}
	}
	broadcastQueue = &memberlist.TransmitLimitedQueue{
		NumNodes: func() int {
			return m.NumMembers()
		},
		RetransmitMult: 3,
	}
	node := m.LocalNode()
	log.Printf("Local node info: %s\n", node.String())
	return nil
}
