package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"nomad_coin/database"
	"nomad_coin/utils"
	"sync"
)

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
	Height   int    `json:"height"`
}

type blockchain struct {
	LastHash string
	Height   int
}

var ErrNotFound = errors.New("block not found")

var db *database.Database

var b *blockchain

var onceForBlockchain sync.Once
var onceForDatabase sync.Once

const DATABASE_FILE_NAME = "blockchain.boltdb"
const BLOCKCHAIN_INFO_BUCKET_NAME = "blockchain"
const BLOCKCHAIN_INFO_KEY_NAME = "checkpoint"

const BLOCK_DATA_BUCKET_NAME = "blockdata"

func GetBlockchainDB() *database.Database {
	if db == nil {
		onceForDatabase.Do(func() {
			db = &database.Database{}
			utils.ErrHandler(db.OpenDB(DATABASE_FILE_NAME))
			utils.ErrHandler(db.CreateBucketWithStringName(BLOCKCHAIN_INFO_BUCKET_NAME))
			utils.ErrHandler(db.CreateBucketWithStringName(BLOCK_DATA_BUCKET_NAME))
		})
	}
	return db
}

func (block *Block) calculateHash() {
	hash := sha256.Sum256([]byte(block.Data + block.PrevHash))
	block.Hash = fmt.Sprintf("%x", hash)
}
func updateBlockchain(newBlock *Block) {
	b.Height = newBlock.Height
	b.LastHash = newBlock.Hash
	byteBlockchainDataToSave, err := utils.ObjectToBytes(b)
	utils.ErrHandler(err)
	utils.ErrHandler(GetBlockchainDB().WriteByteDataToBucket(BLOCKCHAIN_INFO_BUCKET_NAME, BLOCKCHAIN_INFO_KEY_NAME, byteBlockchainDataToSave))
}
func saveNewBlock(newBlock *Block) {
	byteBlockDataToSave, err := utils.ObjectToBytes(&newBlock)
	utils.ErrHandler(err)
	utils.ErrHandler(GetBlockchainDB().WriteByteDataToBucket(BLOCK_DATA_BUCKET_NAME, newBlock.Hash, byteBlockDataToSave))
}
func (b *blockchain) CreateBlockAndSave(data string) {
	newBlock := Block{data, "", b.LastHash, b.Height + 1}
	newBlock.calculateHash()

	updateBlockchain(&newBlock)
	saveNewBlock(&newBlock)
}

func GetBlockByHash(hash string) (*Block, error) {
	data, err := GetBlockchainDB().ReadByteDataFromBucket(BLOCK_DATA_BUCKET_NAME, hash)
	utils.ErrHandler(err)
	if data == nil {
		return nil, ErrNotFound
	}
	block := Block{}
	utils.ErrHandler(utils.ObjectFromBytes(&block, data))
	return &block, nil
}

func (b *blockchain) LoadBlockchain() {
	data, err := GetBlockchainDB().ReadByteDataFromBucket(BLOCKCHAIN_INFO_BUCKET_NAME, BLOCKCHAIN_INFO_KEY_NAME)
	utils.ErrHandler(err)
	if data != nil {
		utils.ErrHandler(utils.ObjectFromBytes(b, data))
	}
}

func GetBlockchain() *blockchain {
	if b == nil {
		onceForBlockchain.Do(func() {

			b = &blockchain{"", 0}
			b.LoadBlockchain()

			if b.Height == 0 {
				b.CreateBlockAndSave("Genesis")
			}
		})
	}
	return b
}

func (b *blockchain) AllBlocks() []*Block {
	allBlocks := []*Block{}
	lastHash := b.LastHash
	for {
		if lastHash == "" {
			break
		}
		block, err := GetBlockByHash(lastHash)
		utils.ErrHandler(err)
		allBlocks = append(allBlocks, block)
		lastHash = block.PrevHash
	}

	for i, j := 0, len(allBlocks)-1; i < j; i, j = i+1, j-1 {
		allBlocks[i], allBlocks[j] = allBlocks[j], allBlocks[i]
	}

	return allBlocks
}
