package blockchain

import "sync"

type memoryPool struct {
	Txs   []*Tx
	mutex sync.Mutex
}

var onceForMempool sync.Once
var mempool *memoryPool

func GetMempool() *memoryPool {
	onceForMempool.Do(func() {
		mempool = &memoryPool{}
	})
	return mempool
}

func (m *memoryPool) clearMempool() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.Txs = []*Tx{}
}

func (m *memoryPool) AddTxToMempool(tx *Tx) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.Txs = append(m.Txs, tx)
}

func isInMempool(txId string) bool {
	GetMempool().mutex.Lock()
	defer GetMempool().mutex.Unlock()

	for _, tx := range GetMempool().Txs {
		for _, txIn := range tx.TxIns {
			if txId == txIn.TxId {
				return true
			}
		}
	}
	return false
}

func GetMempoolTx() []*Tx {
	GetMempool().mutex.Lock()
	defer GetMempool().mutex.Unlock()

	return GetMempool().Txs
}
