package main

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
	Counter       uint64 // Nonce
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.prevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}
