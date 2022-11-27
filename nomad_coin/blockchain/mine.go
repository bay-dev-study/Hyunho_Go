package blockchain

import (
	"nomad_coin/utils"
	"strings"
	"time"
)

func mineNewBlock(block *Block) {
	for {
		targetPrefix := strings.Repeat("0", block.Difficulty)
		block.Timestamp = int(time.Now().Unix())
		hash := utils.HashObject(block)
		// fmt.Println(block)
		// fmt.Println(hash)
		// fmt.Printf("\n")
		if strings.HasPrefix(hash, targetPrefix) {
			block.Hash = hash
			break
		}
		block.Nonce++
	}
}
