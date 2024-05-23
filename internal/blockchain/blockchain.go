package blockchain

import (
	"bytes"
	"encoding/hex"
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

// FindUnspentTransactions returns a list of transactions containing unspent outputs for an address
func (bc *Blockchain) FindUnspentTransactions(pubKeyHash []byte) []transactions.Transaction {
	var unspentTxs []transactions.Transaction
	spentTXOs := make(map[string][]int)

	for _, block := range bc.Blocks {
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				// Check if the output was already spent
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}

				if out.IsLockedWithKey(pubKeyHash) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}

			// Collect spent outputs
			if !tx.IsCoinbase() {
				for _, in := range tx.Vin {
					if in.UsesKey(pubKeyHash) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}
	}

	return unspentTxs
}

// FindUTXO returns unspent transaction outputs for a given address
func (bc *Blockchain) FindUTXO(pubKeyHash []byte) []transactions.TransactionOutput {
	var UTXOs []transactions.TransactionOutput
	unspentTransactions := bc.FindUnspentTransactions(pubKeyHash)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}
