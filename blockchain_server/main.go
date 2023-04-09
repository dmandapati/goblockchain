package main

import (
	"flag"
	"fmt"
	"log"
)

func init() {
	log.SetPrefix("BlockChain_Server: ")
}
func main() {
	port := flag.Uint("port", 3000, "TCP port number for Blockchain server")
	flag.Parse()
	fmt.Println("BlockChain Server 0.0.0.0:", *port)
	app := NewBlockchainServer(uint16(*port))
	// app := NewBlockchainServer(uint16(*port))
	app.Run()
}
