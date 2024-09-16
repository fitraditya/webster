# Webster
*Websocket cluster*

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

A scalable light-weight websocket cluster using gossip protocol.

## Features
* Scalability
  * Each node uses [gossip protocol](https://en.wikipedia.org/wiki/Gossip_protocol) to discover another nodes. New nodes need only be started with the hostname of a single active node in the cluster.
  * When a node recieves an message, it propagates it to the all nodes.

## Installation
Each node can be run as a single binary, docker image, or kubernetes deployment.

### Installing from source
```bash
# download the source code
$ go get github.com/fitraditya/webster

# install the binary
$ go install github.com/fitraditya/webster

# start a node
$ webster server
```

## Run webster

### Flags
```
--ws-address, -a    : websocket server address          default 0.0.0.0:4000
--node-port, -p     : node port for gossip              default 2500
--node-join, -j     : join to existing gossip nodes
--node-config, -c   : node config configuration         default local
```

### Running websocket cluster
```
# open new terminal tab
$ go run main.go server -a 0.0.0.0:4001 -p 2501

# open another new terminal tab
$ go run main.go server -a 0.0.0.0:4002 -p 2502 -j 127.0.0.1:2501
```

In your web browser, open `http://localhost:4001` and `http://localhost:4002` in different tabs. Send message from first tab, and the message will be received by second tab also.

### Running websocket cluster on localhost, LAN, or WAN
Webster using [memberlist](https://github.com/hashicorp/memberlist) as gossip based membership and it has different configurations based on network: default (local), LAN, and WAN. You can specify the configuration when running the websocket server.

```
$ go run main.go server -a 0.0.0.0:4001 -p 2501 -c lan
```

If you run the websocket cluster in different public server, you can use `wan` option.

## To Do
- [ ] Authentication
- [ ] Channel based message subscription
- [ ] Docker and kubernetes deployment


