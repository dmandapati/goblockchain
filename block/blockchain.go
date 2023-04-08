package block

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"goblockchain/utils"
	"log"
	"strings"
	"time"
)

// set the difficulty level, minier sender address and rewards amount for mining
const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 1.0
)

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
	transactionpool   []*Transaction
	chain             []*Block
	blockchainAddress string
}

// function to create Blockchain
func NewBlockchain(blockchainAddress string) *BlockChain {
	b := &Block{}
	bc := new(BlockChain)
	bc.blockchainAddress = blockchainAddress
	bc.CreateBlock(0, b.Hash())
	return bc
}

// function to genrate new block
func (bc *BlockChain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionpool)
	bc.chain = append(bc.chain, b)
	bc.transactionpool = []*Transaction{}
	log.Println("action=New Block, status=success")
	return b
}

// function to identify the last block
func (bc *BlockChain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

// function to print BlockChain
func (bc *BlockChain) Print() {
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
func (bc *BlockChain) AddTransaction(sender string, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {

	t := NewTransaction(sender, recipient, value)

	if sender == MINING_SENDER {
		bc.transactionpool = append(bc.transactionpool, t)
		return true
	}
	if bc.VerifyTransactionSignature(senderPublicKey, s, t) {
		// if bc.CalculateTotalAmount(sender) < value {
		// 	log.Println("ERROR: Not enough balanace in a wallet")
		// 	return false
		// }
		bc.transactionpool = append(bc.transactionpool, t)
		return true
	} else {
		log.Println("ERROR: Verify Transaction")
	}
	return false
}

// function to verify the transaction
func (bc *BlockChain) VerifyTransactionSignature(senderPublicKey *ecdsa.PublicKey, s *utils.Signature, t *Transaction) bool {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
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
	//	log.Println(guessHashStr)
	//	log.Println("Task: mining, action=difficulty check, status=success")
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
	log.Println("action=nonce, status=success")
	return nonce
}

// function for mining
func (bc *BlockChain) Mining() bool {
	bc.AddTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD, nil, nil)
	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)
	log.Println("action=mining, status=success")
	return true
}

// funcation to get the total balance of address
func (bc *BlockChain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value
			if blockchainAddress == t.recipientBlockchainAddress {
				totalAmount += value
			}
			if blockchainAddress == t.senderBlockchainAddress {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}

// create function for creating transactions
func NewTransaction(sender string, recipient string, value float32) *Transaction {
	log.Println("action=transaction, status=success")
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
