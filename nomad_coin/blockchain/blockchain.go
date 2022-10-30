package blockchain

import (
	"errors"
	"nomad_coin/database"
	"nomad_coin/utils"
	"sync"
	"time"
)

type Block struct {
	Hash         string `json:"hash"`
	PrevHash     string `json:"prevHash,omitempty"`
	Height       int    `json:"height"`
	Difficulty   int    `json:"difficulty"`
	Timestamp    int    `json:"timestamp"`
	Nonce        int    `json:"nonce"`
	Transactions []*Tx  `json:"transactions"`
}

type blockchain struct {
	LastHash   string
	Height     int
	Difficulty int
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

const DEFAULT_REWARD_FOR_MINING = 50

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

func (b *blockchain) updateBlockchain(newBlock *Block) {
	b.Height = newBlock.Height
	b.LastHash = newBlock.Hash
	byteBlockchainDataToSave, err := utils.ObjectToBytes(b)
	utils.ErrHandler(err)
	utils.ErrHandler(GetBlockchainDB().WriteByteDataToBucket(BLOCKCHAIN_INFO_BUCKET_NAME, BLOCKCHAIN_INFO_KEY_NAME, byteBlockchainDataToSave))

	if b.Height%RECALCULATE_DIFFICULTY_INTERVAl == 0 {
		b.Difficulty = b.recalculateDifficulty()
	}
}
func (b *blockchain) saveNewBlock(newBlock *Block) {
	byteBlockDataToSave, err := utils.ObjectToBytes(&newBlock)
	utils.ErrHandler(err)
	utils.ErrHandler(GetBlockchainDB().WriteByteDataToBucket(BLOCK_DATA_BUCKET_NAME, newBlock.Hash, byteBlockDataToSave))
}
func (b *blockchain) ConfirmBlock() {
	txSlice := append(GetMempool().Txs, b.makeCoinbaseTx("Hyunho", DEFAULT_REWARD_FOR_MINING))
	newBlock := Block{Hash: "", PrevHash: b.LastHash, Height: b.Height + 1, Difficulty: b.Difficulty, Timestamp: int(time.Now().Unix()), Nonce: 0, Transactions: txSlice}
	newBlock.mine()
	b.saveNewBlock(&newBlock)
	b.updateBlockchain(&newBlock)
	GetMempool().cleanMempool()
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
			b = &blockchain{"", 0, DEFAULT_DIFFICULTY}
			b.LoadBlockchain()

			if b.Height == 0 {
				b.ConfirmBlock()
			}
		})
	}
	return b
}

func (b *blockchain) getBlocksFromLastBlock(number int) []*Block {
	blockSlice := []*Block{}
	lastHash := b.LastHash
	for i := 0; i < number; i++ {
		block, err := GetBlockByHash(lastHash)
		utils.ErrHandler(err)
		blockSlice = append(blockSlice, block)
		lastHash = block.PrevHash
	}
	return blockSlice
}
func (b *blockchain) AllBlocks() []*Block {
	return GetBlockchain().getBlocksFromLastBlock(b.Height)
}
