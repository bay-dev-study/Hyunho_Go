package main

import (
	"nomad_coin/explorer"
	"nomad_coin/rest"
)

func main() {
	go explorer.Start(8080)
	rest.Start(4000)
}
