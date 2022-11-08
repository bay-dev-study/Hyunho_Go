package p2p

import (
	"fmt"
	"nomad_coin/utils"

	"github.com/gorilla/websocket"
)

var Peers map[string]*Peer = make(map[string]*Peer)

type Peer struct {
	key     string
	address string
	port    string
	conn    *websocket.Conn
	inbox   chan []byte
}

func (p *Peer) close() {
	p.conn.Close()
}
func (p *Peer) read() {
	defer p.close()
	for {
		message := &Message{}
		p.conn.ReadJSON(&message)
		fmt.Println("received", message)
		handleMessage(message)
	}
}

func (p *Peer) send() {
	defer p.close()
	for {
		payload, ok := <-p.inbox
		if !ok {
			break
		}
		p.conn.WriteMessage(websocket.TextMessage, payload)
	}
}
func InitPeer(conn *websocket.Conn, address, port string) *Peer {
	key := fmt.Sprintf("%s:%s", address, port)
	peer := &Peer{
		key:     key,
		address: address,
		port:    port,
		conn:    conn,
		inbox:   make(chan []byte),
	}
	Peers[key] = peer
	go peer.read()
	go peer.send()
	sendNewestBlock(peer)
	return peer
}

func AddPeer(address, port, openPort string) {
	uri := fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort[1:])
	conn, _, err := websocket.DefaultDialer.Dial(uri, nil)
	utils.ErrHandler(err)
	InitPeer(conn, address, port)
}
