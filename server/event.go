package server

import (
	"github.com/hashicorp/memberlist"
	"log"
)

type event struct{}

func (e *event) NotifyJoin(node *memberlist.Node) {
	log.Println("A new node has joined: " + node.String())
}

func (e *event) NotifyLeave(node *memberlist.Node) {
	log.Println("A node has left: " + node.String())
}

func (e *event) NotifyUpdate(node *memberlist.Node) {
	log.Println("A node was updated: " + node.String())
}
