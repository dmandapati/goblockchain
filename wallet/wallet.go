package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"goblockchain/utils"

	"github.com/btcsuite/btcutil/base58"

	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	privateKey        *ecdsa.PrivateKey
	publicKey         *ecdsa.PublicKey
	blockchainAddress string
}

// func to marchal json for public/private key and blockchain address
func (w *Wallet) MarshalJSON()([]byte,error){
	return json.Marshal(struct{
		PrivateKey string `json:"private_key"`
		PublicKey string `json:"public_key"`
		BlockchainAddress string `json:"blockchain_address"`
	}{
		PrivateKey: w.PrivateKeyStr(),
		PublicKey: w.PublicKeyStr(),
		BlockchainAddress: w.BlockchainAddress(),

	})
} 


// create bitcoin blockchain address install below go package locally
// go get golang.org/x/crypto/ripemd160
// go get github.com/btcsuite/btcutil/base58

// function to create new wallet
func NewWallet() *Wallet {
	// Step1: Create ECDSA private key (32 bytes) puvlic ket (64 bytes)
	w := new(Wallet)
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	w.privateKey = privateKey
	w.publicKey = &w.privateKey.PublicKey
	// Step2: Performed SHA-256 hashing on public key (32 bytes)
	h2 := sha256.New()
	h2.Write(w.publicKey.X.Bytes())
	h2.Write(w.publicKey.Y.Bytes())
	digest2 := h2.Sum(nil)

	// Step3: Performed RIPEMD-160 hashing on result of SHA-256 (20 bytes)
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)

	// Step4: Add version byte in fornt of RIPEMD-160 hash (0x00 for mainnet)
	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], digest3[:])

	// Step5: Perform SHA-256 hash on the extended RIPEMD-160 result
	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)

	// Step6: Perform SHA-256 hash on result of the previous SHA-256 hash
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)

	// Step7: Take the first 4 bytes of the second SHA-256 hash for checksum
	chsum := digest6[:4]

	// Step8: Add the 4 checksum byes from 7 at the end of extended RIPEMD-160 hash from 4 (25 bytes)
	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4[:])
	copy(dc8[21:], chsum[:])

	// Step9: convert result from a byte string into base58
	address := base58.Encode(dc8)
	w.blockchainAddress = address

	return w
}

// returning privatekey
func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

// output PrivateKey as string
func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

// returning publicKey
func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

// output publicKey as string
func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%x%x", w.publicKey.X.Bytes(), w.publicKey.X.Bytes())
}

// Function to create blockchain address
func (w *Wallet) BlockchainAddress() string {
	return w.blockchainAddress
}




// Create struct for Transaction
type Transaction struct {
	senderPrivateKey           *ecdsa.PrivateKey
	senderPublicKey            *ecdsa.PublicKey
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value                      float32
}

// function for returning new Transaction
func NewTransaction(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, sender string, recipient string, value float32) *Transaction {
	return &Transaction{privateKey, publicKey, sender, recipient, value}
}

// function to genrate tx Signature
func (t *Transaction) GenerateSignature() *utils.Signature {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	r, s, _ := ecdsa.Sign(rand.Reader, t.senderPrivateKey, h[:])
	return &utils.Signature{r, s}
}

// func to convert transaction to json format
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