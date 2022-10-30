package blockchain

import (
	"fmt"
	"nomad_coin/utils"
	"strings"
	"time"
)

const RECALCULATE_DIFFICULTY_INTERVAl int = 5
const TARGET_TIME_INTERVAL_DIFFICULTY int = 10
const TARGET_TIME_INTERVAL_DIFFICULTY_ALLOWANCE int = 3
const DEFAULT_DIFFICULTY int = 2

func (b *blockchain) recalculateDifficulty() int {
	blocks := b.getBlocksFromLastBlock(RECALCULATE_DIFFICULTY_INTERVAl)
	currentTimeInterval := blocks[RECALCULATE_DIFFICULTY_INTERVAl-1].Timestamp/60 - blocks[0].Timestamp/60
	if currentTimeInterval >= TARGET_TIME_INTERVAL_DIFFICULTY+TARGET_TIME_INTERVAL_DIFFICULTY_ALLOWANCE {
		return b.Difficulty - 1
	}
	if currentTimeInterval <= TARGET_TIME_INTERVAL_DIFFICULTY-TARGET_TIME_INTERVAL_DIFFICULTY_ALLOWANCE {
		return b.Difficulty + 1
	}
	return b.Difficulty
}

func (block *Block) mine() {
	for {
		targetPrefix := strings.Repeat("0", block.Difficulty)
		block.Timestamp = int(time.Now().Unix())
		hash := utils.HashObject(block)
		fmt.Println(block)
		fmt.Println(hash)
		fmt.Printf("\n")
		if strings.HasPrefix(hash, targetPrefix) {
			block.Hash = hash
			break
		}
		block.Nonce++
	}
}
