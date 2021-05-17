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
- /set - set value
- /get - get value

Query params expected are `key` and `val`

```shell
# set
curl "http://localhost:8001/set?key=foo&val=bar"
# or use post method 
# curl -d "key=foo&val=bar" http://localhost:8001
# output:ok

# get
curl "http://localhost:8002/get?key=foo"
# output:bar

```

Storage
---

m2 use [badger-db](http://github.com/dgraph-io/badger) as storage

License
---

m2 is under the MIT license. See the LICENSE file for details.