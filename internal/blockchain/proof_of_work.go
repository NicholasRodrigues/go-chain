package main

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
	"strconv"
)

const targetBits = 24

var maxNonce = math.MaxInt64

type ProofOfWork struct {
	block *Block
	T     big.Int
}

// Prepare pow data, this function hashes the data and returns it

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.prevBlockHash,
			pow.block.Data,
			[]byte(strconv.FormatInt(pow.block.Timestamp, 10)),
			[]byte(strconv.FormatInt(targetBits, 10)),
			[]byte(strconv.FormatInt(int64(nonce), 10)),
		},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWork) RunProofOfWork() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	if pow.block == nil {
		pow.block.Data = pow.block.GenesisBlock().Data
	} else {
		temp_block := pow.block
		temp_block.SetHash()
		pow.block = temp_block
	}

	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)

		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(&pow.T) == -1 {
			break
		} else {
			nonce++
		}
	}

	return nonce, hash[:]

}
