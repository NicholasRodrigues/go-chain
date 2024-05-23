package blockchain

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"reflect"
	"time"
)

type Receive func() string
type Input func() string

// ValidateBlockPredicate validates a block using Proof of Work
func ValidateBlockPredicate(pow *ProofOfWork) bool {
	var hashInt big.Int
	data := pow.prepareData(pow.block.Counter)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(&pow.T) == -1
}

// ContentValidatePredicate validates the content of the blockchain
func ContentValidatePredicate(chain *Blockchain) bool {
	if len(chain.Blocks) == 0 {
		return false
	}

	for i := 1; i < len(chain.Blocks); i++ {
		if reflect.DeepEqual(chain.Blocks[i].Hash, chain.Blocks[i-1].Hash) {
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

	newBlock := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(concatData),
		PrevBlockHash: chain.Blocks[len(chain.Blocks)-1].Hash,
		Counter:       round,
	}

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
		data += string(block.Data)
	}
	return data
}

// ChainValidationPredicate validates the blockchain
func ChainValidationPredicate(chain *Blockchain) bool {
	if !ContentValidatePredicate(chain) {
		return false
	}

	for len(chain.Blocks) > 0 {
		lastBlock := chain.Blocks[len(chain.Blocks)-1]
		hash := sha256.Sum256(lastBlock.Data)
		proof := &ProofOfWork{block: lastBlock, T: *big.NewInt(1)}

		if !ValidateBlockPredicate(proof) || !reflect.DeepEqual(lastBlock.Hash, hash) {
			return false
		}

		// Prepare the hash for the next iteration
		hash = [32]byte{}
		copy(hash[:], lastBlock.Data)
		chain.Blocks = chain.Blocks[:len(chain.Blocks)-1]
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
