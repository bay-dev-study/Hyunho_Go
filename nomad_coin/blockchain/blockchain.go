package blockchain

import (
	"errors"
	"fmt"
	"sync"
)

type block struct {
	data     []byte
	hash     [32]byte
	prevHash [32]byte
}

type blockchain struct {
	blocks []block
}

var ErrEmptyBlock = errors.New("empty block error")

var blockchainSingletonInstance *blockchain
var once sync.Once

func GetBlockChain() *blockchain {
	if blockchainSingletonInstance == nil {
		once.Do(func() {
			genesisBlockData := []byte("GenesisBlock")
			genesisBlock := block{genesisBlockData, Hashing(genesisBlockData), [32]byte{}}
			blockchainSingletonInstance = &blockchain{blocks: []block{genesisBlock}}
		})
	}
	return blockchainSingletonInstance
}

func (b *blockchain) getBlockLength() int {
	return len(b.blocks)
}

func (b *blockchain) getLastBlock() *block {
	chainLength := b.getBlockLength()
	return &b.blocks[chainLength-1]
}

func (b *blockchain) AddBlock(data []byte) {
	lastBlock := b.getLastBlock()
	newBlock := block{
		data:     data,
		hash:     Hashing(append(data, lastBlock.prevHash[:]...)),
		prevHash: lastBlock.hash,
	}
	b.blocks = append(b.blocks, newBlock)
}

func (b *blockchain) AddBlockStrData(data string) {
	b.AddBlock([]byte(data))
}

func (b *blockchain) PrintBlocks() {
	blockLength := b.getBlockLength()
	if blockLength == 0 {
		fmt.Println("empty blocks")
		return
	}
	for idx, block := range b.blocks {
		fmt.Printf("[%d Block]\n", idx+1)
		fmt.Printf("Data: %s\n", block.data)
		fmt.Printf("Hash: %x\n", block.hash)
		fmt.Printf("PrevHash: %x\n", block.prevHash)
		fmt.Printf("--------------------\n")
	}
}
