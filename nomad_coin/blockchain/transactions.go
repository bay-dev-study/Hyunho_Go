package blockchain

import (
	"errors"
	"nomad_coin/utils"
	"nomad_coin/wallet"
)

var ErrNotEnoughBalance = errors.New("not enough balance")

var ErrInvalidSignature = errors.New("invalid signature")

type Tx struct {
	TxId      string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

type TxIn struct {
	TxId      string `json:"txId"`
	Index     int    `json:"index"`
	Signature string `json:"signature"`
}

type TxOut struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
}

type UTxOut struct {
	TxId   string `json:"txId"`
	Index  int    `json:"index"`
	Amount int    `json:"amount"`
}

type usedTxOutsFlag map[string]bool

func (tx *Tx) makeTxID() {
	tx.TxId = utils.HashObject(tx)
}

func (tx *Tx) signTxIns() {
	for _, txIn := range tx.TxIns {
		txIn.Signature = wallet.Sign(tx.TxId, wallet.GetWallet())
	}
}

func findTxWithTxId(TxId string) *Tx {
	for _, tx := range getAllTx() {
		if tx.TxId == TxId {
			return tx
		}
	}
	return nil
}

func getAllTx() []*Tx {
	tx_slice := []*Tx{}
	for _, block := range AllBlocks() {
		tx_slice = append(tx_slice, block.Transactions...)
	}
	return tx_slice
}

func makeCoinbaseTx(to string, amount int) *Tx {
	txIn := TxIn{"", -1, "COINBASE"}
	txOut := TxOut{to, amount}
	tx := Tx{TxId: "", Timestamp: utils.GetNowUnixTimestamp(), TxIns: []*TxIn{&txIn}, TxOuts: []*TxOut{&txOut}}
	tx.makeTxID()
	return &tx
}

func checkUsedTxOuts(address string) usedTxOutsFlag {
	allBlocks := AllBlocks()
	usedTxInsMap := usedTxOutsFlag{}
	for _, block := range allBlocks {
		for _, tx := range block.Transactions {
			for _, txIn := range tx.TxIns {
				if txIn.Signature == "COINBASE" {
					break
				}
				txOut := findTxWithTxId(txIn.TxId).TxOuts[txIn.Index]
				if address == txOut.Address {
					usedTxInsMap[txIn.TxId] = true
				}
			}
		}
	}
	return usedTxInsMap
}

func validateTx(tx *Tx) bool {
	for _, txIn := range tx.TxIns {
		txMatchesTxId := findTxWithTxId(txIn.TxId)
		if txMatchesTxId == nil {
			return false
		}
		txOut := txMatchesTxId.TxOuts[txIn.Index]
		isValid := wallet.Verify(txIn.Signature, tx.TxId, txOut.Address)
		if !isValid {
			return false
		}
	}
	return true
}

func GetUTxOfAddress(address string) []*UTxOut {
	allBlocks := AllBlocks()
	usedTxOutsFlag := checkUsedTxOuts(address)
	UTxOutSlice := []*UTxOut{}
	for _, block := range allBlocks {
		for _, tx := range block.Transactions {
			for index, txOut := range tx.TxOuts {
				if txOut.Address == address && !isInMempool(tx.TxId) {
					if _, exists := usedTxOutsFlag[tx.TxId]; !exists {
						uTxOut := UTxOut{
							TxId:   tx.TxId,
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
			TxId:      uTx.TxId,
			Index:     uTx.Index,
			Signature: "",
		})
		if totalAmount >= amount {
			break
		}
	}

	txOuts := []*TxOut{}
	txOuts = append(txOuts, &TxOut{Address: to, Amount: amount})
	change := totalAmount - amount
	if change > 0 {
		txOuts = append(txOuts, &TxOut{Address: from, Amount: change})
	}
	if change < 0 {
		return nil, ErrNotEnoughBalance
	}
	tx := &Tx{
		TxId:      "",
		Timestamp: utils.GetNowUnixTimestamp(),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.makeTxID()
	tx.signTxIns()
	isValid := validateTx(tx)
	if !isValid {
		return nil, ErrInvalidSignature
	}
	return tx, nil
}
