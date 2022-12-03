package blockchain

import (
	"fmt"
	"nomad_coin/utils"
	"strings"
	"time"
)

func mineNewBlock(block *Block) {
	for {
		targetPrefix := strings.Repeat("0", block.Difficulty)
		block.Timestamp = int(time.Now().Unix())
		hash := utils.HashObject(block)

		fmt.Println(hash)
		time.Sleep(1 * time.Millisecond)
		fmt.Printf("\033[1A\033[K") // clear current line
		if strings.HasPrefix(hash, targetPrefix) {
			block.Hash = hash
			break
		}
		block.Nonce++
	}
}
