package blockchain

import (
	"errors"
	"fmt"
	"nomad_coin/database"
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

var ErrNotFound = errors.New("block not found")

var db *database.Database

var b *blockchain = &blockchain{Difficulty: DEFAULT_DIFFICULTY}

var onceForBlockchain sync.Once
var onceForDatabase sync.Once

var databaseFileName string

const DATABASE_FILE_FORMAT = "%s.boltdb"
const BLOCKCHAIN_INFO_BUCKET_NAME = "blockchain"
const BLOCKCHAIN_INFO_KEY_NAME = "checkpoint"
const BLOCK_DATA_BUCKET_NAME = "blockdata"

const DEFAULT_REWARD_FOR_MINING = 50

func (b *blockchain) updateBlockchain(newBlock *Block) {
	b.Height = newBlock.Height
	b.LastHash = newBlock.Hash

	persistBlockhain(b)
	if b.Height%RECALCULATE_DIFFICULTY_INTERVAl == 0 {
		b.Difficulty = recalculateDifficulty()
	}
}

func (b *blockchain) ConfirmBlock() {
	txSlice := append(GetMempool().Txs, makeCoinbaseTx(wallet.GetWallet().Address, DEFAULT_REWARD_FOR_MINING))
	newBlock := Block{Hash: "", PrevHash: b.LastHash, Height: b.Height + 1, Difficulty: b.Difficulty, Timestamp: int(time.Now().Unix()), Nonce: 0, Transactions: txSlice}
	newBlock.mine()
	saveNewBlock(&newBlock)
	b.updateBlockchain(&newBlock)
	GetMempool().cleanMempool()
}

func (b *blockchain) ReplaceAllBlocks(allBlocks []*Block) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.Difficulty = allBlocks[0].Difficulty
	b.Height = allBlocks[0].Height
	b.LastHash = allBlocks[0].Hash
	persistBlockhain(b)
	db.EmptyBucket(BLOCK_DATA_BUCKET_NAME)
	for _, block := range allBlocks {
		saveNewBlock(block)
	}
}

func (b *blockchain) AddPeerBlock(newBlock *Block) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	mempool.mutex.Lock()
	defer mempool.mutex.Unlock()

	b.Height = newBlock.Height
	b.Difficulty = newBlock.Difficulty
	b.LastHash = newBlock.Hash

	persistBlockhain(b)
	saveNewBlock(newBlock)

	GetMempool().cleanMempool()
}

func persistBlockhain(b *blockchain) {
	byteBlockchainDataToSave, err := utils.ObjectToBytes(b)
	utils.ErrHandler(err)
	utils.ErrHandler(GetBlockchainDB().WriteByteDataToBucket(BLOCKCHAIN_INFO_BUCKET_NAME, BLOCKCHAIN_INFO_KEY_NAME, byteBlockchainDataToSave))
}

func saveNewBlock(newBlock *Block) {
	byteBlockDataToSave, err := utils.ObjectToBytes(&newBlock)
	utils.ErrHandler(err)
	utils.ErrHandler(GetBlockchainDB().WriteByteDataToBucket(BLOCK_DATA_BUCKET_NAME, newBlock.Hash, byteBlockDataToSave))
}

func getBlocksFromLastBlock(number int) []*Block {

	blockSlice := []*Block{}
	lastHash := GetBlockchain().LastHash
	for i := 0; i < number; i++ {
		block, err := GetBlockByHash(lastHash)
		utils.ErrHandler(err)
		blockSlice = append(blockSlice, block)
		lastHash = block.PrevHash
	}
	return blockSlice
}

func getAllTx() []*Tx {
	tx_slice := []*Tx{}
	for _, block := range AllBlocks() {
		tx_slice = append(tx_slice, block.Transactions...)
	}
	return tx_slice
}
func findTxWithTxId(TxId string) *Tx {
	for _, tx := range getAllTx() {
		if tx.TxId == TxId {
			return tx
		}
	}
	return nil
}

func SetBlockchainDatabaseFileName(port string) {
	databaseFileName = fmt.Sprintf(DATABASE_FILE_FORMAT, port)
}

func GetBlockchainDB() *database.Database {
	onceForDatabase.Do(func() {
		db = &database.Database{}
		utils.ErrHandler(db.OpenDB(databaseFileName))
		utils.ErrHandler(db.CreateBucketWithStringName(BLOCKCHAIN_INFO_BUCKET_NAME))
		utils.ErrHandler(db.CreateBucketWithStringName(BLOCK_DATA_BUCKET_NAME))
	})
	return db
}

func GetNewestBlock() *Block {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	lastHash := GetBlockchain().LastHash
	block, err := GetBlockByHash(lastHash)
	utils.ErrHandler(err)
	return block
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

func LoadBlockchain() {
	data, err := GetBlockchainDB().ReadByteDataFromBucket(BLOCKCHAIN_INFO_BUCKET_NAME, BLOCKCHAIN_INFO_KEY_NAME)
	utils.ErrHandler(err)
	if data != nil {
		utils.ErrHandler(utils.ObjectFromBytes(b, data))
	}
}

func GetBlockchain() *blockchain {
	onceForBlockchain.Do(func() {
		LoadBlockchain()

		if b.Height == 0 {
			b.ConfirmBlock()
		}
	})
	return b
}

func AllBlocks() []*Block {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return getBlocksFromLastBlock(GetBlockchain().Height)
}
