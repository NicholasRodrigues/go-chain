package main

import (
	"bytes"
	"math/big"
	"strconv"
)

const targetBits = 24

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
