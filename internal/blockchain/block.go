package blockchain

import (
	"bytes"
	"crypto/sha256"
	"github.com/NicholasRodrigues/go-chain/internal/transactions"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64
	Transactions  []*transactions.Transaction
	PrevBlockHash []byte
	Hash          []byte
	Counter       int // Nonce
}

// SetHash recalculates the hash of the block.
func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.HashTransactions(), timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

// HashTransactions returns a hash of the transactions in the block.
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.Serialize())
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

// NewBlock creates and returns a new block.
func NewBlock(transactions []*transactions.Transaction, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Counter:       0,
	}

	pow := NewProofOfWork(block)
	hash, nonce := pow.Run()
	block.Hash = hash[:]
	block.Counter = nonce

	return block
}
