package main

import (
	"flag"
	"fmt"
	"github.com/qichengzx/m2/server"
	"log"
	"net/http"
	"strings"
)

var (
	members = flag.String("members", "", "list of members")
	port    = flag.Int("port", 8001, "http port")
	db      = flag.String("db", "syncmap", "db type")
	dir     = flag.String("dir", "data", "db dir")
)

func main() {
	flag.Parse()
	var memberList []string
	if *members != "" {
		memberList = strings.Split(*members, ",")
	}

	if *db == "rocksdb" && *dir == "" {
		log.Fatalln("dir is needed when using rocksdb mode")
	}

	server := server.New(*db)
	if err := server.Start(*port, memberList); err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/set", server.SetHandler)
	http.HandleFunc("/del", server.DelHandler)
	http.HandleFunc("/get", server.GetHandler)
	fmt.Println("Listening on:", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Fatalln(err)
	}
}
