package blockchain

import (
	"bytes"
	"github.com/NicholasRodrigues/go-chain/internal/transactions"
)

type Blockchain struct {
	Blocks []*Block
}

// AddBlock adds a new block to the blockchain with the given transactions.
func (bc *Blockchain) AddBlock(transactions []*transactions.Transaction) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(transactions, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

// NewGenesisBlock creates and returns the genesis block.
func NewGenesisBlock() *Block {
	coinbase := transactions.NewTransaction(
		[]transactions.TransactionInput{
			{Txid: []byte{}, Vout: -1, ScriptSig: "Genesis"},
		},
		[]transactions.TransactionOutput{
			{Value: 50, ScriptPubKey: "coinbase"},
		},
	)
	return NewBlock([]*transactions.Transaction{coinbase}, []byte{})
}

// NewBlockchain creates and returns a new blockchain with the genesis block.
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

// IsValid validates the blockchain.
func (bc *Blockchain) IsValid() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		prevBlock := bc.Blocks[i-1]

		if !bytes.Equal(currentBlock.PrevBlockHash, prevBlock.Hash) {
			return false
		}

		pow := NewProofOfWork(currentBlock)
		if !pow.Validate() {
			return false
		}
	}

	return true
}
