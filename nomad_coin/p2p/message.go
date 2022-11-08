package p2p

import (
	"encoding/json"
	"nomad_coin/blockchain"
	"nomad_coin/utils"
)

type MessageKind int

const (
	MessageNewestBlock MessageKind = iota
	MessageAllBlocksRequest
	MessageAllBlocksResponse
)

type Message struct {
	Kind    MessageKind `json:"kind"`
	Payload []byte      `json:"payload"`
}

func makeBytesMessage(kind MessageKind, payload interface{}) []byte {
	message := &Message{Kind: kind, Payload: utils.ToJson(&payload)}
	return utils.ToJson(&message)
}

func sendNewestBlock(p *Peer) {
	block := blockchain.GetNewestBlock()
	p.inbox <- makeBytesMessage(MessageNewestBlock, &block)
}

func handleMessage(message *Message) {
	switch message.Kind {
	case MessageNewestBlock:
		newestBlock := &blockchain.Block{}
		err := json.Unmarshal(message.Payload, &newestBlock)
		utils.ErrHandler(err)
	}
}
