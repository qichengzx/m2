m2
--
m2 is a simple http key/value cache system based on [hashicorp/memberlist](https://github.com/hashicorp/memberlist).

memberlist is a Go library using a [gossip](https://www.consul.io/docs/architecture/gossip)  based protocol.

so, m2 is a distributed key/value cache system.

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
./m2 --port 8001
```

will show node address and http port:
```shell
2021/01/02 09:52:17 Local node info: 192.168.0.6:18001
Listening on: 8001
```

then, start second node with first node as part of the member list:
```shell
./m2 --port 8002 --members "192.168.0.6:18001"
```

will show:
```shell
2021/01/02 09:57:30 [DEBUG] memberlist: Initiating push/pull sync with:  192.168.0.6:18001
2021/01/02 09:57:30 Local node info: 192.168.0.6:18002
Listening on: 8002

```

The first node will show:
```shell
2021/01/02 09:57:51 [DEBUG] memberlist: Initiating push/pull sync with: your-host-name-8002 192.168.0.6:18002
2021/01/02 09:57:53 [DEBUG] memberlist: Stream connection from=192.168.0.6:55311
```

Key/Value Api
---

HTTP API
- /set - set value
- /get - get value
- /del - delete value

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

# delete
curl "http://localhost:8003/del?key=foo"
# output:ok
```

Storage
---

m2 support sync.Map and [rocksdb](https://github.com/tecbot/gorocksdb) as storage. default is sync.Map.

so data will be lost when server is down.

License
---

m2 is under the MIT license. See the LICENSE file for details.