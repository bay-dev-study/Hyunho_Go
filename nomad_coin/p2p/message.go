package p2p

import (
	"encoding/json"
	"fmt"
	"nomad_coin/blockchain"
	"nomad_coin/utils"
	"strings"
)

type MessageKind int

const (
	MessageNewestBlock MessageKind = iota
	MessageAllBlocksRequest
	MessageAllBlocksResponse
	MessageNewPeerNotify
	MessageNewBlockNotify
	MessageNewTransactionNotify
)

type Message struct {
	Kind    MessageKind `json:"kind"`
	Payload []byte      `json:"payload"`
}

func makeBytesMessage(kind MessageKind, payload interface{}) []byte {
	message := &Message{Kind: kind, Payload: utils.ToJson(&payload)}
	return utils.ToJson(&message)
}

func requestAllBlocks(p *Peer) {
	p.inbox <- makeBytesMessage(MessageAllBlocksRequest, nil)
}

func sendAllBlocks(p *Peer) {
	allBlocks := blockchain.AllBlocks()
	p.inbox <- makeBytesMessage(MessageAllBlocksResponse, &allBlocks)
}

func sendNewestBlock(p *Peer) {
	block := blockchain.GetNewestBlock()
	fmt.Println("Sending newest block", block)
	p.inbox <- makeBytesMessage(MessageNewestBlock, block)
}

func notifyNewPeer(p *Peer, address, port, openPort string) {
	p.inbox <- makeBytesMessage(MessageNewPeerNotify, fmt.Sprintf("%s:%s:%s", address, port, openPort))
}

func notifyNewBlock(p *Peer, newBlock *blockchain.Block) {
	p.inbox <- makeBytesMessage(MessageNewBlockNotify, newBlock)
}

func notifyNewTransaction(p *Peer, tx *blockchain.Tx) {
	p.inbox <- makeBytesMessage(MessageNewTransactionNotify, tx)
}

func handleMessage(message *Message, peer *Peer) {
	switch message.Kind {
	case MessageNewestBlock:
		newestBlock := &blockchain.Block{}
		err := json.Unmarshal(message.Payload, &newestBlock)
		utils.ErrHandler(err)
		recentBlockHeight := blockchain.GetNewestBlock().Height
		if newestBlock.Height > recentBlockHeight {
			fmt.Printf("Requesting all blocks from %s\n", peer.key)
			requestAllBlocks(peer)
		}

	case MessageAllBlocksRequest:
		fmt.Printf("%s wants all the blocks\n", peer.key)
		sendAllBlocks(peer)

	case MessageAllBlocksResponse:
		fmt.Printf("Received all the blocks from %s\n", peer.key)
		var allBlocks []*blockchain.Block
		err := json.Unmarshal(message.Payload, &allBlocks)
		utils.ErrHandler(err)
		blockchain.GetBlockchain().ReplaceAllBlocks(allBlocks)

	case MessageNewPeerNotify:
		var addPeerPayload string
		err := json.Unmarshal(message.Payload, &addPeerPayload)
		utils.ErrHandler(err)
		addPeerPayloadSlice := strings.Split(addPeerPayload, ":")
		AddPeer(addPeerPayloadSlice[0], addPeerPayloadSlice[1], addPeerPayloadSlice[2], false)

	case MessageNewBlockNotify:
		var newBlock *blockchain.Block
		err := json.Unmarshal(message.Payload, &newBlock)
		utils.ErrHandler(err)
		blockchain.GetBlockchain().AddPeerBlock(newBlock)

	case MessageNewTransactionNotify:
		var newTx *blockchain.Tx
		err := json.Unmarshal(message.Payload, &newTx)
		utils.ErrHandler(err)
		blockchain.GetMempool().AddTxToMempool(newTx)

	}
}
