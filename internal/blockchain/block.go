package blockchain

import (
	"bytes"
	"crypto/sha256"
	"strconv"
)

type Block struct {
	Timestamp     int64
	Data          []byte // Transactions
	prevBlockHash []byte
	Hash          []byte
	Counter       int // Nonce
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.prevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     0,
		Data:          []byte(data),
		prevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Counter:       0,
	}

	block.SetHash()
	return block
}
