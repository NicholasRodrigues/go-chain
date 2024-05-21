package main

import (
	"fmt"
	"math/big"
	"reflect"
	"time"
)

type Receive func() string
type Input func() string

type ProofOfWork struct {
	block *Block
	T     big.Int
}

// Simple Content Validation Predicate implementation from backbone protocol
func ContentValidatePredicate(x *Blockchain) bool {

	if len(x.blocks) == 0 {
		return false
	}

	for i := range x.blocks {
		if i == 0 {
			continue
		}
		if reflect.DeepEqual(x.blocks[i].Hash, x.blocks[i-1].Hash) {
			return false
		}
	}
	return true
}

func InputContributionFunction(data []byte, cr *Blockchain, round int, input Input, receive Receive) {

	input_data := input()
	receive_data := receive()

	concat_data := input_data + receive_data
	// creating new block

	newBlock := &Block{time.Now().Unix(), []byte(concat_data), cr.blocks[len(cr.blocks)-1].Hash, []byte{}, round}
	cr.blocks = append(cr.blocks, newBlock)

	if !ContentValidatePredicate(cr) {
		fmt.Println("Content Validation Failed")
		cr.blocks = cr.blocks[:len(cr.blocks)-1]
	} else {
		fmt.Println("Content Validation Passed")
	}
}

// Function to read the chain
func ChainReadFunction(c *Blockchain) string {
	data := ""
	for i := range c.blocks {
		data += string(c.blocks[i].Data)
	}
	return data
}
