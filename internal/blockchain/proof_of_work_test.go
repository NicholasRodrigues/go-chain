package blockchain

import (
	"math/big"
	"testing"
)

func TestNewProofOfWork(t *testing.T) {
	block := NewBlock("test data", []byte{})
	pow := NewProofOfWork(block)

	expectedTarget := big.NewInt(1)
	expectedTarget.Lsh(expectedTarget, uint(256-Difficulty))

	if pow == nil {
		t.Fatal("expected non-nil ProofOfWork")
	}
	if pow.block != block {
		t.Error("expected pow.block to be the same as block")
	}
	if pow.T.Cmp(expectedTarget) != 0 {
		t.Errorf("expected target %s, got %s", expectedTarget, &pow.T)
	}
}
