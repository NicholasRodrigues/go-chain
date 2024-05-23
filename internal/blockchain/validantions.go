package blockchain

import (
	"bytes"
	"fmt"
	"github.com/NicholasRodrigues/go-chain/internal/transactions"
)

// Receive is a function type for receiving data
type Receive func() string

// Input is a function type for providing input data
type Input func() string

// ValidateBlockPredicate validates a block using Proof of Work
func ValidateBlockPredicate(pow *ProofOfWork) bool {
	return pow.Validate()
}

// ContentValidatePredicate validates the content of the blockchain
func ContentValidatePredicate(chain *Blockchain) bool {
	if len(chain.Blocks) == 0 {
		return false
	}

	for i := 1; i < len(chain.Blocks); i++ {
		if !bytes.Equal(chain.Blocks[i].PrevBlockHash, chain.Blocks[i-1].Hash) {
			return false
		}
	}
	return true
}

// InputContributionFunction processes input and adds a new block to the blockchain
func InputContributionFunction(data []byte, chain *Blockchain, round int, input Input, receive Receive) {
	inputData := input()
	receiveData := receive()
	concatData := inputData + receiveData

	newTransaction := transactions.NewTransaction(
		[]transactions.TransactionInput{{Txid: []byte{}, Vout: -1, ScriptSig: concatData}},
		[]transactions.TransactionOutput{{Value: 50, ScriptPubKey: "pubkey1"}},
	)

	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	newBlock := NewBlock([]*transactions.Transaction{newTransaction}, prevBlock.Hash)

	chain.Blocks = append(chain.Blocks, newBlock)

	if !ContentValidatePredicate(chain) {
		fmt.Println("Content Validation Failed")
		chain.Blocks = chain.Blocks[:len(chain.Blocks)-1]
	} else {
		fmt.Println("Content Validation Passed")
	}
}

// ChainReadFunction returns the data of the blockchain as a string
func ChainReadFunction(chain *Blockchain) string {
	var data string
	for _, block := range chain.Blocks {
		for _, tx := range block.Transactions {
			data += fmt.Sprintf("%x", tx.ID)
		}
	}
	return data
}

// ChainValidationPredicate validates the blockchain
func ChainValidationPredicate(chain *Blockchain) bool {
	if !ContentValidatePredicate(chain) {
		return false
	}

	for _, block := range chain.Blocks {
		pow := NewProofOfWork(block)
		if !pow.Validate() {
			return false
		}
	}

	return true
}

// MaxChain finds the best chain among multiple chains
func MaxChain(chains [][]Blockchain) []*Block {
	var bestChain []*Block

	for i := 1; i < len(chains); i++ {
		if ChainValidationPredicate(&chains[i][0]) {
			bestChain = maxBlocks(bestChain, chains[i][0].Blocks)
		}
	}

	return bestChain
}

// maxBlocks returns the longer of two block slices
func maxBlocks(a, b []*Block) []*Block {
	if len(a) > len(b) {
		return a
	}
	return b
}
