package main

import (
	"nomad_coin/wallet"
)

func main() {
	wallet.Start()
	// defer database.CloseAllOpenedDB()
	// cli.Start()
}
