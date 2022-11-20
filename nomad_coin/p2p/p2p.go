package p2p

import (
	"fmt"
	"nomad_coin/blockchain"
	"nomad_coin/utils"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var Peers map[string]*Peer = make(map[string]*Peer)
var peersMutex = sync.Mutex{}

type Peer struct {
	key     string
	address string
	port    string
	conn    *websocket.Conn
	inbox   chan []byte
}

func (p *Peer) close() {
	peersMutex.Lock()
	defer peersMutex.Unlock()
	p.conn.Close()
}
func (p *Peer) read() {
	defer p.close()
	for {
		message := &Message{}
		p.conn.ReadJSON(&message)
		handleMessage(message, p)
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
	peersMutex.Lock()
	defer peersMutex.Unlock()

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
	return peer
}

func AddPeer(address, port, openPort string, isBroadcast bool) {
	uri := fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort)
	conn, _, err := websocket.DefaultDialer.Dial(uri, nil)
	utils.ErrHandler(err)
	peer := InitPeer(conn, address, port)
	if isBroadcast {
		broadcastNewPeer(address, port)
	}
	sendNewestBlock(peer)
}

func broadcastNewPeer(address, port string) {
	for _, peer := range Peers {
		if strings.Compare(peer.address, address) == 0 && strings.Compare(peer.port, port) == 0 {
			continue
		}
		notifyNewPeer(peer, address, port, peer.port)
	}
}

func BroadcastNewBlock(newBlock *blockchain.Block) {
	for _, peer := range Peers {
		notifyNewBlock(peer, newBlock)
	}
}

func BroadcastNewTransaction(tx *blockchain.Tx) {
	for _, peer := range Peers {
		notifyNewTransaction(peer, tx)
	}
}
