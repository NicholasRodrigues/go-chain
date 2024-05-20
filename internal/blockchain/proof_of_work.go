package blockchain

import (
	"math/big"
)

const Difficulty = 16

type ProofOfWork struct {
	block *Block
	T     big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{b, *target}

	return pow
}
