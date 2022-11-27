package p2p

import (
	"encoding/json"
	"fmt"
	"nomad_coin/blockchain"
	"nomad_coin/utils"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

type PeersMap struct {
	Peers map[string]*Peer
	mutex sync.Mutex
}

var PeersMapInstance = PeersMap{Peers: make(map[string]*Peer)}

type Peer struct {
	key     string
	address string
	port    string
	conn    *websocket.Conn
	inbox   chan []byte
}

func WritePeersToJsonEncoder(encoder *json.Encoder) {
	PeersMapInstance.mutex.Lock()
	defer PeersMapInstance.mutex.Unlock()

	encoder.Encode(PeersMapInstance.Peers)
}

func (p *Peer) close() {
	PeersMapInstance.mutex.Lock()
	defer PeersMapInstance.mutex.Unlock()

	delete(PeersMapInstance.Peers, p.key)
	p.conn.Close()
}

func (p *Peer) read() {
	defer p.close()
	for {
		message := &Message{}
		err := p.conn.ReadJSON(&message)
		if err != nil {
			break
		}
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
	PeersMapInstance.mutex.Lock()
	defer PeersMapInstance.mutex.Unlock()

	key := fmt.Sprintf("%s:%s", address, port)
	peer := &Peer{
		key:     key,
		address: address,
		port:    port,
		conn:    conn,
		inbox:   make(chan []byte),
	}
	PeersMapInstance.Peers[key] = peer
	go peer.read()
	go peer.send()

	sendNewestBlock(peer)
	return peer
}

func AddPeer(address, port, openPort string, needToBroadcast bool) {
	uri := fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort)
	conn, _, err := websocket.DefaultDialer.Dial(uri, nil)
	utils.ErrHandler(err)
	InitPeer(conn, address, port)
	if needToBroadcast {
		broadcastNewPeer(address, port)
	}
}

func broadcastNewPeer(address, port string) {
	PeersMapInstance.mutex.Lock()
	defer PeersMapInstance.mutex.Unlock()

	for _, peer := range PeersMapInstance.Peers {
		if strings.Compare(peer.address, address) == 0 && strings.Compare(peer.port, port) == 0 {
			continue
		}
		notifyNewPeer(peer, address, port, peer.port)
	}
}

func BroadcastNewBlock(newBlock *blockchain.Block) {
	PeersMapInstance.mutex.Lock()
	defer PeersMapInstance.mutex.Unlock()

	for _, peer := range PeersMapInstance.Peers {
		notifyNewBlock(peer, newBlock)
	}
}

func BroadcastNewTransaction(tx *blockchain.Tx) {
	PeersMapInstance.mutex.Lock()
	defer PeersMapInstance.mutex.Unlock()

	for _, peer := range PeersMapInstance.Peers {
		notifyNewTransaction(peer, tx)
	}
}
