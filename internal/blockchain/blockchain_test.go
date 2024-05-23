package blockchain

import (
	"bytes"
	"testing"
)

func TestProofOfWork_Run(t *testing.T) {
	block := NewBlock("test data", []byte{})
	pow := NewProofOfWork(block)

	hash, nonce := pow.Run()
	block.Hash = hash[:]
	block.Counter = nonce

	if !pow.Validate() {
		t.Errorf("Proof of work validation failed for block with hash %x", block.Hash)
	}
}

func TestProofOfWork_Validate(t *testing.T) {
	block := NewBlock("test data", []byte{})
	pow := NewProofOfWork(block)

	hash, nonce := pow.Run()
	block.Hash = hash[:]
	block.Counter = nonce

	if !pow.Validate() {
		t.Errorf("Expected block to be valid, but validation failed")
	}
}

func TestNewBlockchain(t *testing.T) {
	bc := NewBlockchain()
	if bc == nil {
		t.Fatalf("Expected blockchain to be created")
	}
	if len(bc.Blocks) != 1 {
		t.Fatalf("Expected blockchain to have one block, got %d", len(bc.Blocks))
	}
	if !bytes.Equal(bc.Blocks[0].Hash, NewGenesisBlock().Hash) {
		t.Errorf("Expected genesis block hash to match")
	}
}

func TestAddBlock(t *testing.T) {
	bc := NewBlockchain()
	bc.AddBlock("First Block after Genesis")
	bc.AddBlock("Second Block after Genesis")

	if len(bc.Blocks) != 3 {
		t.Fatalf("Expected blockchain to have three Blocks, got %d", len(bc.Blocks))
	}
	if !bytes.Equal(bc.Blocks[1].PrevBlockHash, bc.Blocks[0].Hash) {
		t.Errorf("Expected first block to reference genesis block hash")
	}
	if !bytes.Equal(bc.Blocks[2].PrevBlockHash, bc.Blocks[1].Hash) {
		t.Errorf("Expected second block to reference first block hash")
	}
}

func TestBlockchain_IsValid(t *testing.T) {
	bc := NewBlockchain()
	bc.AddBlock("First Block after Genesis")
	bc.AddBlock("Second Block after Genesis")

	if !bc.IsValid() {
		t.Errorf("Expected blockchain to be valid")
	}
}
