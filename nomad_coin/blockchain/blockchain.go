package blockchain

import (
	"errors"
	"fmt"
)

type Block struct {
	data     []byte
	hash     [32]byte
	prevHash [32]byte
}

type BlockChain struct {
	blocks []Block
}

var ErrEmptyBlock = errors.New("empty block error")

func (blockchain *BlockChain) getBlockLength() int {
	return len(blockchain.blocks)
}

func (blockchain *BlockChain) getLastBlock() (*Block, error) {
	chainLength := blockchain.getBlockLength()
	if chainLength > 0 {
		return &blockchain.blocks[chainLength-1], nil
	}
	return nil, ErrEmptyBlock
}

func (blockchain *BlockChain) addBlock(data []byte) {
	lastBlock, err := blockchain.getLastBlock()
	var newBlock Block
	if err != nil {
		newBlock = Block{
			data:     data,
			hash:     Hashing(data),
			prevHash: [32]byte{},
		}
	} else {
		newBlock = Block{
			data:     data,
			hash:     Hashing(append(data, lastBlock.prevHash[:]...)),
			prevHash: lastBlock.hash,
		}
	}
	blockchain.blocks = append(blockchain.blocks, newBlock)
}

func (blockchain *BlockChain) AddBlockStrData(data string) {
	blockchain.addBlock([]byte(data))
}

func (blockchain *BlockChain) PrintBlocks() {
	blockLength := blockchain.getBlockLength()
	if blockLength == 0 {
		fmt.Println("empty blocks")
		return
	}
	for idx, block := range blockchain.blocks {
		fmt.Printf("[%d Block]\n", idx)
		fmt.Printf("--------------------\n")
		fmt.Printf("Data: %s\n", block.data)
		fmt.Printf("Hash: %x\n", block.hash)
		fmt.Printf("PrevHash: %x\n", block.prevHash)
	}
}

// chain := BlockChain{}
// chain.addBlock("Genesis Block")
// chain.addBlock("Second Block")
// chain.addBlock("Third Block")
// chain.listBlocks()
