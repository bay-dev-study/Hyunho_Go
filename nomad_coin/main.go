package main

import (
	"nomad_coin/cli"
	"nomad_coin/database"
)

func main() {
	defer database.CloseAllOpenedDB()
	cli.Start()
}
