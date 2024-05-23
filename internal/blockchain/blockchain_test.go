package blockchain

import (
	"bytes"
	"github.com/NicholasRodrigues/go-chain/internal/transactions"
	"testing"
)

// TestBlockchainAddBlock checks if blocks are added correctly to the blockchain
func TestBlockchainAddBlock(t *testing.T) {
	bc := NewBlockchain()
	tx := createTransaction()
	bc.AddBlock([]*transactions.Transaction{tx})

	if len(bc.Blocks) != 2 {
		t.Errorf("Expected blockchain length 2, but got %d", len(bc.Blocks))
	}

	if !bytes.Equal(bc.Blocks[1].PrevBlockHash, bc.Blocks[0].Hash) {
		t.Errorf("Previous block hash does not match")
	}
}

// TestBlockchainIsValid checks if the blockchain validation works correctly
func TestBlockchainIsValid(t *testing.T) {
	bc := NewBlockchain()
	tx := createTransaction()
	bc.AddBlock([]*transactions.Transaction{tx})
	bc.AddBlock([]*transactions.Transaction{tx})

	if !bc.IsValid() {
		t.Errorf("Blockchain should be valid")
	}

	// Tamper with the blockchain
	bc.Blocks[1].Transactions = []*transactions.Transaction{
		transactions.NewTransaction(
			[]transactions.TransactionInput{{Txid: []byte("tampered"), Vout: 0, ScriptSig: "tampered"}},
			[]transactions.TransactionOutput{{Value: 50, ScriptPubKey: "tampered"}},
		),
	}

	if bc.IsValid() {
		t.Errorf("Blockchain should be invalid after tampering")
	}
}
