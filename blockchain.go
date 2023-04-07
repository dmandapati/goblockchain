package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

// set the difficulty level for mining
const MINING_DIFFICULTY = 3

// Block Struct
type Block struct {
	timestamp    int64
	nonce        int
	previousHash [32]byte
	transactions []*Transaction
}

// funcation for create new block
func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	b := new(Block)
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash
	b.transactions = transactions
	return b
}

// function to print block
func (b *Block) Print() {
	fmt.Printf("timestamp       %d\n", b.timestamp)
	fmt.Printf("nonce           %d\n", b.nonce)
	fmt.Printf("previous_hash   %x\n", b.previousHash)
	for _, t := range b.transactions {
		t.Print()
	}
}

// function to create hash
func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256([]byte(m))
}

// function to marcel the json format
func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
		PreviousHash [32]byte       `json:"previous_hash"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
	})
}

// Blockchain Struct
type BlockChain struct {
	transactionpool []*Transaction
	chain           []*Block
}

// function to create Blockchain
func NewBlockchain() *BlockChain {
	b := &Block{}
	bc := new(BlockChain)
	bc.CreateBlock(0, b.Hash())
	return bc
}

// function to genrate new block
func (bc *BlockChain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionpool)
	bc.chain = append(bc.chain, b)
	bc.transactionpool = []*Transaction{}
	return b
}

// function to identify the last block
func (bc *BlockChain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

// function to print BlockChain
func (bc *BlockChain) print() {
	for i, block := range bc.chain {
		fmt.Printf("%s chain %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

// Creating Transaction struct
type Transaction struct {
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value                      float32
}

// This function will add the transaction to transaction pool
func (bc *BlockChain) addTransaction(sender string, recipient string, value float32) {
	t := NewTransaction(sender, recipient, value)
	bc.transactionpool = append(bc.transactionpool, t)
}

// Function to copy transactions from transactions pool
func (bc *BlockChain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionpool {
		transactions = append(transactions, NewTransaction(
			t.senderBlockchainAddress,
			t.recipientBlockchainAddress,
			t.value))
	}
	return transactions
}

// function to compute the difficulty level or mining
func (bc *BlockChain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{timestamp: 0, nonce: nonce, previousHash: previousHash, transactions: transactions}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	fmt.Println(guessHashStr)
	return guessHashStr[:difficulty] == zeros
}

// function to return nonce value
func (bc *BlockChain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}
	return nonce
}

// create function for creating transactions
func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{sender, recipient, value}
}

// Create function to print the transaction
func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" sender_blockchain_address    %s\n", t.senderBlockchainAddress)
	fmt.Printf(" recipient_blockchain_address %s\n", t.recipientBlockchainAddress)
	fmt.Printf(" value                        %f\n", t.value)
}

// Create function to marshel the tranaction in json format
func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.senderBlockchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Value:     t.value,
	})
}

func init() {
	log.SetPrefix("BlockChain: ")
}

func main() {
	blockChain := NewBlockchain()
	blockChain.print()

	blockChain.addTransaction("A", "B", 1.0)

	previousHash := blockChain.LastBlock().Hash()
	nonce := blockChain.ProofOfWork()
	blockChain.CreateBlock(nonce, previousHash)
	blockChain.print()

	blockChain.addTransaction("C", "D", 2.0)
	blockChain.addTransaction("X", "Y", 3.0)

	previousHash = blockChain.LastBlock().Hash()
	nonce = blockChain.ProofOfWork()
	blockChain.CreateBlock(nonce, previousHash)
	blockChain.print()
}
