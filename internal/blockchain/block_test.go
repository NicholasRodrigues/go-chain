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
		PrevBlockHash: prevHash,
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
	if block.PrevBlockHash == nil || !bytes.Equal(block.PrevBlockHash, prevHash) {
		t.Errorf("expected previous hash %x, got %x", prevHash, block.PrevBlockHash)
	}
	if block.Hash == nil || len(block.Hash) == 0 {
		t.Error("expected non-nil and non-empty hash")
	}
}

func TestNewGenesisBlock(t *testing.T) {
	genesisBlock := NewGenesisBlock()
	if genesisBlock == nil {
		t.Fatalf("Expected genesis block to be created")
	}
	if len(genesisBlock.Hash) == 0 {
		t.Errorf("Expected genesis block hash to be set")
	}
	if !bytes.Equal(genesisBlock.PrevBlockHash, []byte{}) {
		t.Errorf("Expected genesis block to have no previous block hash")
	}
}
