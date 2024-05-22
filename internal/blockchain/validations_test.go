package blockchain

import (
	"testing"
)

func TestContentValidatePredicate(t *testing.T) {
	bc := NewBlockchain()
	bc.AddBlock("New Block 1")
	bc.AddBlock("New Block 2")

	if !ContentValidatePredicate(bc) {
		t.Error("expected blockchain to be valid")
	}

	bc.blocks[1].Hash = bc.blocks[0].Hash
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

	if len(bc.blocks) != 2 {
		t.Errorf("expected blockchain length 2, got %d", len(bc.blocks))
	}

	concat_data := "input data" + "received data"
	if string(bc.blocks[1].Data) != concat_data {
		t.Errorf("expected block data %s, got %s", concat_data, string(bc.blocks[1].Data))
	}

	if !ContentValidatePredicate(bc) {
		t.Error("expected blockchain to be valid after contribution")
	}
}

func TestChainReadFunction(t *testing.T) {
	bc := NewBlockchain()
	bc.AddBlock("New Block 1")
	bc.AddBlock("New Block 2")

	data := ChainReadFunction(bc)
	expectedData := "Genesis BlockNew Block 1New Block 2"

	if data != expectedData {
		t.Errorf("expected chain data %s, got %s", expectedData, data)
	}
}
