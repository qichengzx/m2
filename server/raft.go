package server

import (
	"github.com/hashicorp/raft"
	"net/http"
)

func (h Server) RaftJoin(w http.ResponseWriter, r *http.Request) {
	nodeid := r.FormValue("nodeid")
	addr := r.FormValue("addr")
	if h.raft.State() != raft.Leader {
		w.Write([]byte("error:not leader"))
		return
	}

	configFuture := h.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		w.Write([]byte("error:get raft configuration"))
		return
	}

	f := h.raft.AddVoter(raft.ServerID(nodeid), raft.ServerAddress(addr), 0, 0)
	if f.Error() != nil {
		w.Write([]byte("error:add voter"))
		return
	}

	w.Write([]byte("ok"))
}
