package blockchain

import (
	"errors"
	"nomad_coin/utils"
	"sync"
)

var onceForMempool sync.Once

type Tx struct {
	TxID      string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

func (tx *Tx) makeTxID() {
	tx.TxID = utils.HashObject(tx)
}

type TxIn struct {
	TxID  string `json:"txId"`
	Index int    `json:"index"`
	Owner string `json:"owner"`
}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

type UTxOut struct {
	TxID   string `json:"txId"`
	Index  int    `json:"index"`
	Amount int    `json:"amount"`
}

func makeCoinbaseTx(to string, amount int) *Tx {
	txIn := TxIn{"", -1, "COINBASE"}
	txOut := TxOut{to, amount}
	tx := Tx{TxID: "", Timestamp: utils.GetNowUnixTimestamp(), TxIns: []*TxIn{&txIn}, TxOuts: []*TxOut{&txOut}}
	tx.makeTxID()
	return &tx
}

type usedTxOutsFlag map[string]bool

func checkUsedTxOuts(address string) usedTxOutsFlag {
	allBlocks := AllBlocks()
	usedTxInsMap := usedTxOutsFlag{}
	for _, block := range allBlocks {
		for _, tx := range block.Transactions {
			for _, txIn := range tx.TxIns {
				if address == txIn.Owner {
					usedTxInsMap[txIn.TxID] = true
				}
			}
		}
	}
	return usedTxInsMap
}
func GetUTxOfAddress(address string) []*UTxOut {
	allBlocks := AllBlocks()
	usedTxOutsFlag := checkUsedTxOuts(address)
	UTxOutSlice := []*UTxOut{}
	for _, block := range allBlocks {
		for _, tx := range block.Transactions {
			for index, txOut := range tx.TxOuts {
				if txOut.Owner == address && !isInMempool(tx.TxID) {
					if _, exists := usedTxOutsFlag[tx.TxID]; !exists {
						uTxOut := UTxOut{
							TxID:   tx.TxID,
							Index:  index,
							Amount: txOut.Amount,
						}
						UTxOutSlice = append(UTxOutSlice, &uTxOut)
						break
					}
				}
			}
		}
	}
	return UTxOutSlice
}

func BalanceByAddress(address string) int {
	uTxSlice := GetUTxOfAddress(address)
	totalAmount := 0
	for _, uTx := range uTxSlice {
		totalAmount += uTx.Amount
	}
	return totalAmount
}

func MakeTx(from, to string, amount int) (*Tx, error) {
	uTxSlice := GetUTxOfAddress(from)
	totalAmount := 0
	txIns := []*TxIn{}
	for _, uTx := range uTxSlice {
		totalAmount += uTx.Amount
		txIns = append(txIns, &TxIn{
			TxID:  uTx.TxID,
			Index: uTx.Index,
			Owner: from,
		})
		if totalAmount >= amount {
			break
		}
	}

	txOuts := []*TxOut{}
	txOuts = append(txOuts, &TxOut{Owner: to, Amount: amount})
	change := totalAmount - amount
	if change > 0 {
		txOuts = append(txOuts, &TxOut{Owner: from, Amount: change})
	}
	if change < 0 {
		return nil, errors.New("not enough balance")
	}
	tx := Tx{
		TxID:      "",
		Timestamp: utils.GetNowUnixTimestamp(),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.makeTxID()
	return &tx, nil
}

type memoryPool struct {
	Txs []*Tx
}

var mempool *memoryPool

func (m *memoryPool) cleanMempool() {
	m.Txs = []*Tx{}
}

func (m *memoryPool) AddTx(from, to string, amount int) error {
	tx, err := MakeTx(from, to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func isInMempool(txId string) bool {
	for _, tx := range GetMempool().Txs {
		for _, txIn := range tx.TxIns {
			if txId == txIn.TxID {
				return true
			}
		}
	}
	return false
}

func GetMempool() *memoryPool {
	onceForMempool.Do(func() {
		mempool = &memoryPool{}
	})
	return mempool
}
