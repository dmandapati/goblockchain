package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// Block Struct
type Block struct {
	nonce        int
	previousHash string
	timestamp    int64
	transactions []string
}

// funcation for create new block
func NewBlock(nonce int, previousHash string) *Block {
	b := new(Block)
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash
	return b
}

// function to print block
func (b *Block) Print() {
	fmt.Printf("timestamp       %d\n", b.timestamp)
	fmt.Printf("nonce           %d\n", b.nonce)
	fmt.Printf("previous_hash   %s\n", b.previousHash)
	fmt.Printf("transactions    %s\n", b.transactions)
}

// Blockchain Struct
type BlockChain struct {
	transactionpool []string
	chain           []*Block
}

// function to create Blockchain
func NewBlockchain() *BlockChain {
	bc := new(BlockChain)
	bc.CreateBlock(0, "init hash")
	return bc
}

// function to genrate new block
func (bc *BlockChain) CreateBlock(nonce int, previousHash string) *Block {
	b := NewBlock(nonce, previousHash)
	bc.chain = append(bc.chain, b)
	return b
}

// function to print BlockChain
func (bc *BlockChain) print() {
	for i, block := range bc.chain {
		fmt.Printf("%s chain %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

func init() {
	log.SetPrefix("BlockChain: ")
}

func main() {
	blockChain := NewBlockchain()
	blockChain.print()
	blockChain.CreateBlock(5, "hash 1")
	blockChain.print()
	blockChain.CreateBlock(2, "hash 2")
	blockChain.print()
}