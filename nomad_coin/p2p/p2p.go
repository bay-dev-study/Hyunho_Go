package p2p

import (
	"fmt"
	"nomad_coin/utils"

	"github.com/gorilla/websocket"
)

var Peers map[string]*Peer = make(map[string]*Peer)

type Peer struct {
	conn *websocket.Conn
}

func InitPeer(conn *websocket.Conn, address, port string) {
	p := &Peer{
		conn,
	}
	key := fmt.Sprintf("%s:%s", address, port)
	Peers[key] = p
}

func AddPeer(address, port, openPort string) {
	uri := fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort[1:])
	conn, _, err := websocket.DefaultDialer.Dial(uri, nil)
	utils.ErrHandler(err)
	InitPeer(conn, address, port)
}
