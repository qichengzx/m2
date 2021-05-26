package server

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/raft"
	"net/http"
)

func (s Server) RaftJoin(w http.ResponseWriter, r *http.Request) {
	nodeid := r.FormValue("nodeid")
	addr := r.FormValue("addr")
	if s.raft.State() != raft.Leader {
		w.Write([]byte("error:not leader"))
		return
	}

	configFuture := s.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		w.Write([]byte("error:get raft configuration"))
		return
	}

	f := s.raft.AddVoter(raft.ServerID(nodeid), raft.ServerAddress(addr), 0, 0)
	if f.Error() != nil {
		w.Write([]byte("error:add voter"))
		return
	}

	w.Write([]byte("ok"))
}

func (s Server) RaftStatus(w http.ResponseWriter, r *http.Request) {
	j, err := json.Marshal(s.raft.Stats())
	if err != nil {
		w.Write([]byte("error:marshal raft status"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func (s Server) RaftLeave(w http.ResponseWriter, r *http.Request) {
	nodeid := r.FormValue("nodeid")
	if s.raft.State() != raft.Leader {
		w.Write([]byte("error:not the leader"))
		return
	}

	configFuture := s.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		w.Write([]byte("error:get raft configuration"))
		return
	}

	future := s.raft.RemoveServer(raft.ServerID(nodeid), 0, 0)
	if err := future.Error(); err != nil {
		w.Write([]byte(fmt.Sprintf("error:remove node %s: %s", nodeid, err.Error())))
		return
	}

	w.Write([]byte("removed successfully"))
}
