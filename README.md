m2
--
m2 is a simple http key/value cache system based on [hashicorp/raft](https://github.com/hashicorp/raft).

---

Install
---

```shell
go get github.com/qichengzx/m2
```

Create Cluster
---

Start first node
```shell
./m2 --node_id 1 --port 8001 --raft_port 18001
```

then, start second node with first node as part of the member list:
```shell
./m2 --node_id 2 --port 8002 --raft_port 18002
```

```join cluster
curl -d 'nodeid=2&addr=127.0.0.1:18002' http://localhost:8001/raft/join
```

Key/Value Api
---

HTTP API
- /set - set key&value
- /get - get value
- /del - del key

Query params expected are `key` and `val`

```shell
# set
curl "http://localhost:8001/set?key=foo&val=bar"
# or use post method 
# curl -d "key=foo&val=bar" http://localhost:8001/set
# output:ok

# get
curl "http://localhost:8002/get?key=foo"
# output:bar

# del
curl "http://localhost:8001/del?key=foo"
# output:ok
```

Raft Api
---

HTTP API
- /raft/join - join raft cluster
- /raft/leave - leave raft cluster
- /raft/status - get raft node status

```shell
# join
curl "http://localhost:8001/raft/join?nodeid=2&addr=127.0.0.1:18002"
# or use post method 
# curl -d "nodeid=2&addr=127.0.0.1:18002" http://localhost:8001/raft/join
# output:ok

# leave
curl "http://localhost:8001/raft/leave?nodeid=2&addr=127.0.0.1:18002"
# output:removed successfully

# node status
curl "http://localhost:8001/raft/status"
# output:
{
    "applied_index": "2",
    "commit_index": "2",
    "fsm_pending": "0",
    "last_contact": "0",
    "last_log_index": "2",
    "last_log_term": "2",
    "last_snapshot_index": "0",
    "last_snapshot_term": "0",
    "latest_configuration": "[{Suffrage:Voter ID:1 Address:127.0.0.1:18001}]",
    "latest_configuration_index": "0",
    "num_peers": "0",
    "protocol_version": "3",
    "protocol_version_max": "3",
    "protocol_version_min": "0",
    "snapshot_version_max": "1",
    "snapshot_version_min": "0",
    "state": "Leader",
    "term": "2"
}
```

Storage
---

m2 use [badger-db](http://github.com/dgraph-io/badger) as storage

License
---

m2 is under the MIT license. See the LICENSE file for details.