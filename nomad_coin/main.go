package main

import "nomad_coin/blockchain"

func main() {
	chain := blockchain.GetBlockChain()
	chain.AddBlockStrData("Second Block")
	chain.AddBlockStrData("Third Block")
	chain.PrintBlocks()
}
