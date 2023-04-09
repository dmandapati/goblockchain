package main

import (
	"goblockchain/block"
	"goblockchain/wallet"
	"io"
	"log"
	"net/http"
	"strconv"
)

var cache map[string]*block.BlockChain = make(map[string]*block.BlockChain)

type BlockchainServer struct {
	port uint16
}

func NewBlockchainServer(port uint16) *BlockchainServer {
	return &BlockchainServer{port}
}

func (bsc *BlockchainServer) Port() uint16 {
	return bsc.port
}

func (bsc *BlockchainServer) GetBlockchain() *block.BlockChain {
	bc, ok := cache["blockchain"]
	if !ok {
		minersWallet := wallet.NewWallet()
		bc = block.NewBlockchain(minersWallet.BlockchainAddress(), bsc.Port())
		cache["blockchain"] = bc
		log.Printf(("private_key %v"), minersWallet.PrivateKeyStr())
		log.Printf(("public_key %v"), minersWallet.PublicKeyStr())
		log.Printf(("blockchain_address %v"), minersWallet.BlockchainAddress())
	}
	return bc
}

func (bsc *BlockchainServer) GetChain(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bsc.GetBlockchain()
		m, _ := bc.MarshalJSON()
		io.WriteString(w, string(m[:]))
	default:
		log.Printf("ERROR: Invalid HTTP Method")
	}
}

func (bsc *BlockchainServer) Run() {
	http.HandleFunc("/", bsc.GetChain)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(bsc.port)), nil))
}
