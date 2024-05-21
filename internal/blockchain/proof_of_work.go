package main

import (
	"math/big"
)

type ProofOfWork struct {
	block *Block
	T     big.Int
}
