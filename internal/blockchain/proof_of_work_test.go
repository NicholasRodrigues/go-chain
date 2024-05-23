package blockchain

import (
	"github.com/NicholasRodrigues/go-chain/internal/transactions"
	"testing"
)

// TestProofOfWorkRun ensures that the proof of work runs correctly
func TestProofOfWorkRun(t *testing.T) {
	tx := createTransaction()
	block := NewBlock([]*transactions.Transaction{tx}, []byte{})
	pow := NewProofOfWork(block)
	hash, nonce := pow.Run()

	block.Hash = hash[:]
	block.Counter = nonce

	if !pow.Validate() {
		t.Errorf("Proof of work did not validate")
	}
}

// TestProofOfWorkValidate checks if proof of work validation works correctly
func TestProofOfWorkValidate(t *testing.T) {
	tx := createTransaction()
	block := NewBlock([]*transactions.Transaction{tx}, []byte{})
	pow := NewProofOfWork(block)

	if !pow.Validate() {
		t.Errorf("Proof of work validation failed")
	}
}
