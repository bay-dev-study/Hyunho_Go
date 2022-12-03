package blockchain

import (
	"fmt"
	"nomad_coin/database"
	"nomad_coin/utils"
	"sync"
)

type blockchainDB struct {
	db    *database.Database
	mutex sync.Mutex
}

var dbInstance *blockchainDB

var databaseFileName string

const DATABASE_FILE_FORMAT = "%s.boltdb"
const BLOCKCHAIN_INFO_BUCKET_NAME = "blockchain"
const BLOCKCHAIN_INFO_KEY_NAME = "checkpoint"
const BLOCK_DATA_BUCKET_NAME = "blockdata"

var onceForDatabase sync.Once

func getBlockchainDB() *blockchainDB {
	onceForDatabase.Do(func() {
		dbInstance = &blockchainDB{}
		dbInstance.db = &database.Database{}
		utils.ErrHandler(dbInstance.db.OpenDB(databaseFileName))
		utils.ErrHandler(dbInstance.db.CreateBucketWithStringName(BLOCKCHAIN_INFO_BUCKET_NAME))
		utils.ErrHandler(dbInstance.db.CreateBucketWithStringName(BLOCK_DATA_BUCKET_NAME))
	})
	return dbInstance
}

func saveNewBlockToDB(newBlock *Block) {
	getBlockchainDB().mutex.Lock()
	defer getBlockchainDB().mutex.Unlock()
	// defer fmt.Println("block saved")

	byteBlockDataToSave, err := utils.ObjectToBytes(&newBlock)
	utils.ErrHandler(err)
	utils.ErrHandler(getBlockchainDB().db.WriteByteDataToBucket(BLOCK_DATA_BUCKET_NAME, newBlock.Hash, byteBlockDataToSave))
}

func saveBlockchainToDB(chain *blockchain) {
	getBlockchainDB().mutex.Lock()
	defer getBlockchainDB().mutex.Unlock()

	byteBlockchainDataToSave, err := utils.ObjectToBytes(chain)
	utils.ErrHandler(err)
	utils.ErrHandler(getBlockchainDB().db.WriteByteDataToBucket(BLOCKCHAIN_INFO_BUCKET_NAME, BLOCKCHAIN_INFO_KEY_NAME, byteBlockchainDataToSave))
}

func loadBlockchainFromDatabase() *blockchain {
	getBlockchainDB().mutex.Lock()
	defer getBlockchainDB().mutex.Unlock()

	data, err := getBlockchainDB().db.ReadByteDataFromBucket(BLOCKCHAIN_INFO_BUCKET_NAME, BLOCKCHAIN_INFO_KEY_NAME)
	utils.ErrHandler(err)
	chain := &blockchain{Difficulty: DEFAULT_DIFFICULTY}
	if data != nil {
		utils.ErrHandler(utils.ObjectFromBytes(&chain, data))
	}
	return chain
}

func clearBlockDB() {
	getBlockchainDB().mutex.Lock()
	defer getBlockchainDB().mutex.Unlock()

	getBlockchainDB().db.EmptyBucket(BLOCK_DATA_BUCKET_NAME)
}

func SetBlockchainDatabaseFileName(port string) {
	databaseFileName = fmt.Sprintf(DATABASE_FILE_FORMAT, port)
}

func GetBlockByHash(hash string) (*Block, error) {
	getBlockchainDB().mutex.Lock()
	defer getBlockchainDB().mutex.Unlock()

	data, err := getBlockchainDB().db.ReadByteDataFromBucket(BLOCK_DATA_BUCKET_NAME, hash)
	utils.ErrHandler(err)
	if data == nil {
		return nil, ErrNotFound
	}
	block := Block{}
	utils.ErrHandler(utils.ObjectFromBytes(&block, data))
	return &block, nil
}
