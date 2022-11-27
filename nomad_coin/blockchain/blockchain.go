package blockchain

import (
	"encoding/json"
	"errors"
	"nomad_coin/utils"
	"nomad_coin/wallet"
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
	mutex      sync.Mutex
}

const RECALCULATE_DIFFICULTY_INTERVAl int = 5
const TARGET_TIME_INTERVAL_DIFFICULTY int = 10
const TARGET_TIME_INTERVAL_DIFFICULTY_ALLOWANCE int = 3
const DEFAULT_DIFFICULTY int = 2

var ErrNotFound = errors.New("block not found")

var blockchainInstance *blockchain

const DEFAULT_REWARD_FOR_MINING = 50

var onceForBlockchain sync.Once

func GetBlockchain() *blockchain {
	onceForBlockchain.Do(func() {
		blockchainInstance = loadBlockchainFromDatabase()

		if blockchainInstance.Height == 0 {
			blockchainInstance.CreateNewBlockFromTx()
		}
	})
	return blockchainInstance
}

func (chain *blockchain) updateBlockchain(newBlock *Block) {
	chain.Height = newBlock.Height
	chain.LastHash = newBlock.Hash
	chain.Difficulty = newBlock.Difficulty

	saveBlockchainToDB(chain)
}

func (b *blockchain) ReplaceAllBlocks(blocksToReplace []*Block) {
	GetBlockchain().mutex.Lock()
	defer GetBlockchain().mutex.Unlock()

	b.updateBlockchain(blocksToReplace[0])
	clearBlockDB()
	for _, block := range blocksToReplace {
		saveNewBlockToDB(block)
	}
}

func (b *blockchain) AddPeerBlock(newBlock *Block) {
	GetBlockchain().mutex.Lock()
	defer GetBlockchain().mutex.Unlock()

	GetBlockchain().updateBlockchain(newBlock)
	saveNewBlockToDB(newBlock)
	GetMempool().clearMempool()
}

func calculateBlockchainDifficulty(currentHeight, currentDifficulty int) int {

	if currentHeight%RECALCULATE_DIFFICULTY_INTERVAl == 0 && currentHeight != 0 {
		blocks := getBlocksFromLastBlock(RECALCULATE_DIFFICULTY_INTERVAl)
		currentTimeInterval := blocks[RECALCULATE_DIFFICULTY_INTERVAl-1].Timestamp/60 - blocks[0].Timestamp/60
		if currentTimeInterval >= TARGET_TIME_INTERVAL_DIFFICULTY+TARGET_TIME_INTERVAL_DIFFICULTY_ALLOWANCE {
			return currentDifficulty - 1
		}
		if currentTimeInterval <= TARGET_TIME_INTERVAL_DIFFICULTY-TARGET_TIME_INTERVAL_DIFFICULTY_ALLOWANCE {
			return currentDifficulty + 1
		}
	}
	return currentDifficulty
}

func getBlocksFromLastBlock(howMany int) []*Block {
	blockSlice := []*Block{}
	lastHash := GetBlockchain().LastHash
	for i := 0; i < howMany; i++ {
		block, err := GetBlockByHash(lastHash)
		utils.ErrHandler(err)
		blockSlice = append(blockSlice, block)
		lastHash = block.PrevHash
	}
	return blockSlice
}

func GetNewestBlock() *Block {
	GetBlockchain().mutex.Lock()
	defer GetBlockchain().mutex.Unlock()

	lastHash := GetBlockchain().LastHash
	block, err := GetBlockByHash(lastHash)
	utils.ErrHandler(err)
	return block
}

func (chain *blockchain) CreateNewBlockFromTx() {
	chain.mutex.Lock()
	defer chain.mutex.Unlock()

	txSlice := append(GetMempoolTx(), makeCoinbaseTx(wallet.GetWallet().Address, DEFAULT_REWARD_FOR_MINING))
	newBlock := Block{Hash: "", PrevHash: chain.LastHash, Height: chain.Height + 1, Difficulty: calculateBlockchainDifficulty(chain.Height, chain.Difficulty), Timestamp: int(time.Now().Unix()), Nonce: 0, Transactions: txSlice}
	mineNewBlock(&newBlock)
	saveNewBlockToDB(&newBlock)
	GetMempool().clearMempool()
	chain.updateBlockchain(&newBlock)
}

func AllBlocks() []*Block {
	GetBlockchain().mutex.Lock()
	defer GetBlockchain().mutex.Unlock()

	return getBlocksFromLastBlock(GetBlockchain().Height)
}

func WriteBlockchainToJsonEncoder(encoder *json.Encoder) {
	GetBlockchain().mutex.Lock()
	defer GetBlockchain().mutex.Unlock()

	encoder.Encode(GetBlockchain())
}
