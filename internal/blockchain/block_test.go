package blockchain

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"testing"
	"time"
)

func TestSetHash(t *testing.T) {
	data := "test data"
	prevHash := []byte("previous hash")
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		prevBlockHash: prevHash,
		Counter:       0,
	}
	block.SetHash()

	expectedHash := sha256.Sum256(bytes.Join([][]byte{prevHash, []byte(data), []byte(strconv.FormatInt(block.Timestamp, 10))}, []byte{}))
	if !bytes.Equal(block.Hash, expectedHash[:]) {
		t.Errorf("expected hash %x, got %x", expectedHash, block.Hash)
	}
}

func TestNewBlock(t *testing.T) {
	data := "test data"
	prevHash := []byte("previous hash")
	block := NewBlock(data, prevHash)

	if block.Data == nil || !bytes.Equal(block.Data, []byte(data)) {
		t.Errorf("expected data %s, got %s", data, string(block.Data))
	}
	if block.prevBlockHash == nil || !bytes.Equal(block.prevBlockHash, prevHash) {
		t.Errorf("expected previous hash %x, got %x", prevHash, block.prevBlockHash)
	}
	if block.Hash == nil || len(block.Hash) == 0 {
		t.Error("expected non-nil and non-empty hash")
	}
}

func TestNewGenesisBlock(t *testing.T) {
	block := NewGenesisBlock()

	if block.Data == nil || string(block.Data) != "Genesis Block" {
		t.Errorf("expected data 'Genesis Block', got %s", string(block.Data))
	}
	if block.prevBlockHash == nil || len(block.prevBlockHash) != 0 {
		t.Error("expected empty previous hash")
	}
	if block.Hash == nil || len(block.Hash) == 0 {
		t.Error("expected non-nil and non-empty hash")
	}
}

func TestNewBlockchain(t *testing.T) {
	bc := NewBlockchain()

	if len(bc.blocks) != 1 {
		t.Errorf("expected blockchain length 1, got %d", len(bc.blocks))
	}
	if string(bc.blocks[0].Data) != "Genesis Block" {
		t.Errorf("expected genesis block data 'Genesis Block', got %s", string(bc.blocks[0].Data))
	}
}
