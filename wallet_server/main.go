package main

import (
	"flag"
	"log"
)

func init() {
	log.SetPrefix("Wallet_Server: ")
}
func main() {
	port := flag.Uint("port", 8080, "TCP port number for Wallet server")
	gateway := flag.String("gateway", "http://127.0.0.1:3000", "Blockchain Gateway")
	flag.Parse()

	app := NewWalletServer(uint16(*port), *gateway)
	app.Run()
}
