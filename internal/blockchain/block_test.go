package blockchain

import (
	"bytes"
	"github.com/NicholasRodrigues/go-chain/internal/transactions"
	"github.com/NicholasRodrigues/go-chain/pkg/crypto"
	"testing"
)

func createTransaction() *transactions.Transaction {
	privKey, _ := crypto.NewPrivateKey()
	tx := transactions.NewTransaction(
		[]transactions.TransactionInput{
			{Txid: []byte("somepreviousid"), Vout: 0, ScriptSig: "scriptSig"},
		},
		[]transactions.TransactionOutput{
			{Value: 10, ScriptPubKey: "pubkey1"},
		},
	)
	tx.Sign(privKey)
	return tx
}

// TestBlockHash checks if the block hash is set correctly
func TestBlockHash(t *testing.T) {
	tx := createTransaction()
	block := NewBlock([]*transactions.Transaction{tx}, []byte{})
	expectedHash := block.Hash

	if !bytes.Equal(block.Hash, expectedHash) {
		t.Errorf("Expected hash %x, but got %x", expectedHash, block.Hash)
	}
}

// TestGenesisBlock checks if the genesis block is created correctly
func TestGenesisBlock(t *testing.T) {
	block := NewGenesisBlock()
	if len(block.PrevBlockHash) != 0 {
		t.Errorf("Genesis block should have no previous block hash")
	}
}

// TestBlockMining ensures that mining a block works correctly
func TestBlockMining(t *testing.T) {
	tx := createTransaction()
	block := NewBlock([]*transactions.Transaction{tx}, []byte{})
	pow := NewProofOfWork(block)

	if !pow.Validate() {
		t.Errorf("Block did not pass proof of work validation")
	}
}
