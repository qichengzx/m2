package main

import (
	"flag"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
	"github.com/qichengzx/m2/fsm"
	"github.com/qichengzx/m2/server"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	tcpTimeout = 1 * time.Second
	snapInterval = 30 * time.Second
	snapThreshold = 1000
)

var (
	nodeID   = flag.String("node_id", "node_1", "raft node id")
	port     = flag.Int("port", 8001, "http port")
	raftaddr = flag.String("raft_addr", "18001", "raft addr")
	dir      = flag.String("store_dir", "data", "db dir")
)

func main() {
	flag.Parse()

	raftConf := raft.DefaultConfig()
	raftConf.LocalID = raft.ServerID(*nodeID)
	raftConf.SnapshotInterval = snapInterval
	raftConf.SnapshotThreshold = snapThreshold

	badgerOpt := badger.DefaultOptions(*dir)
	badgerDB, err := badger.Open(badgerOpt)
	if err != nil {
		log.Fatal(err)
		return
	}

	fsmStore := fsm.NewBadger(badgerDB)

	store, err := raftboltdb.NewBoltStore(filepath.Join(*dir, "raft"))
	if err != nil {
		log.Fatal(err)
		return
	}

	cacheStore, err := raft.NewLogCache(256, store)
	if err != nil {
		log.Fatal(err)
		return
	}

	snapshotStore, err := raft.NewFileSnapshotStore(*dir, 1, os.Stdout)
	if err != nil {
		log.Fatal(err)
		return
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", *raftaddr)
	if err != nil {
		log.Fatal(err)
		return
	}
	transport, err := raft.NewTCPTransport(*raftaddr, tcpAddr, 3, tcpTimeout, os.Stdout)
	if err != nil {
		log.Fatal(err)
		return
	}

	raftServer, err := raft.NewRaft(raftConf, fsmStore, cacheStore, store, snapshotStore, transport)
	if err != nil {
		log.Fatal(err)
		return
	}

	raftServer.BootstrapCluster(raft.Configuration{
		Servers: []raft.Server{
			{
				ID:      raft.ServerID(*nodeID),
				Address: transport.LocalAddr(),
			},
		},
	})

	server := server.New(raftServer, badgerDB)
	http.HandleFunc("/raft/join", server.RaftJoin)
	http.HandleFunc("/raft/status", server.RaftStatus)
	http.HandleFunc("/raft/leave", server.RaftLeave)
	http.HandleFunc("/set", server.SetHandler)
	http.HandleFunc("/get", server.GetHandler)
	http.HandleFunc("/del", server.DelHandler)
	fmt.Println("HTTP Server Listening on:", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Fatalln(err)
	}
}
