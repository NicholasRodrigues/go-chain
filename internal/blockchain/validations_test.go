package blockchain

import (
	"fmt"
	"github.com/NicholasRodrigues/go-chain/internal/transactions"
	"testing"
)

func TestContentValidatePredicate(t *testing.T) {
	bc := NewBlockchain()
	tx1 := createTransaction()
	tx2 := createTransaction()
	bc.AddBlock([]*transactions.Transaction{tx1})
	bc.AddBlock([]*transactions.Transaction{tx2})

	if !ContentValidatePredicate(bc) {
		t.Error("expected blockchain to be valid")
	}

	// Tamper with the blockchain
	bc.Blocks[1].Transactions[0].Vin[0].ScriptSig = "tampered"
	bc.Blocks[1].SetHash() // Recalculate the hash after tampering
	if ContentValidatePredicate(bc) {
		t.Error("expected blockchain to be invalid")
	}
}

func TestInputContributionFunction(t *testing.T) {
	bc := NewBlockchain()
	round := 1
	data := []byte("initial data")

	input := func() string {
		return "input data"
	}
	receive := func() string {
		return "received data"
	}

	InputContributionFunction(data, bc, round, input, receive)

	if len(bc.Blocks) != 2 {
		t.Errorf("expected blockchain length 2, got %d", len(bc.Blocks))
	}

	concatData := "input data" + "received data"
	if string(bc.Blocks[1].Transactions[0].Vin[0].ScriptSig) != concatData {
		t.Errorf("expected block data %s, got %s", concatData, string(bc.Blocks[1].Transactions[0].Vin[0].ScriptSig))
	}

	if !ContentValidatePredicate(bc) {
		t.Error("expected blockchain to be valid after contribution")
	}
}

func TestChainReadFunction(t *testing.T) {
	bc := NewBlockchain()
	tx1 := createTransaction()
	tx2 := createTransaction()
	bc.AddBlock([]*transactions.Transaction{tx1})
	bc.AddBlock([]*transactions.Transaction{tx2})

	data := ChainReadFunction(bc)
	expectedData := fmt.Sprintf("%x%x%x", bc.Blocks[0].Transactions[0].ID, tx1.ID, tx2.ID)

	if data != expectedData {
		t.Errorf("expected chain data %s, got %s", expectedData, data)
	}
}
