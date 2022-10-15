package main

import "nomad_coin/blockchain"

func main() {
	chain := blockchain.BlockChain{}
	chain.AddBlockStrData("Genesis Block")
	chain.AddBlockStrData("Second Block")
	chain.AddBlockStrData("Third Block")
	chain.PrintBlocks()
}
