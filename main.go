package main

import (
	"flag"
	"fmt"
	"github.com/dgraph-io/badger"
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
)

var (
	nodeID   = flag.String("node_id", "node_1", "raft node id")
	port     = flag.Int("port", 8001, "http port")
	raftport = flag.String("raft_port", "18001", "raft port")
	dir      = flag.String("store_dir", "data", "db dir")
)

func main() {
	flag.Parse()

	raftConf := raft.DefaultConfig()
	raftConf.LocalID = raft.ServerID(*nodeID)
	raftConf.SnapshotThreshold = 1024

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

	raftAddr := fmt.Sprintf("127.0.0.1:%s", *raftport)
	tcpAddr, err := net.ResolveTCPAddr("tcp", raftAddr)
	if err != nil {
		log.Fatal(err)
		return
	}
	transport, err := raft.NewTCPTransport(raftAddr, tcpAddr, 3, tcpTimeout, os.Stdout)
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
	http.HandleFunc("/set", server.SetHandler)
	http.HandleFunc("/get", server.GetHandler)
	http.HandleFunc("/del", server.DelHandler)
	fmt.Println("HTTP Server Listening on:", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Fatalln(err)
	}
}
