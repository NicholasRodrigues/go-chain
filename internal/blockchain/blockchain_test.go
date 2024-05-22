package blockchain

import (
	"bytes"
	"testing"
)

func TestAddBlock(t *testing.T) {
	bc := NewBlockchain()
	bc.AddBlock("New Block 1")
	bc.AddBlock("New Block 2")

	if len(bc.blocks) != 3 {
		t.Errorf("expected blockchain length 3, got %d", len(bc.blocks))
	}
	if string(bc.blocks[1].Data) != "New Block 1" {
		t.Errorf("expected block 1 data 'New Block 1', got %s", string(bc.blocks[1].Data))
	}
	if string(bc.blocks[2].Data) != "New Block 2" {
		t.Errorf("expected block 2 data 'New Block 2', got %s", string(bc.blocks[2].Data))
	}
	if !bytes.Equal(bc.blocks[2].prevBlockHash, bc.blocks[1].Hash) {
		t.Errorf("expected block 2 previous hash to match block 1 hash")
	}
}
