package blockchain

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/NicholasRodrigues/go-chain/internal/transactions"
	"github.com/NicholasRodrigues/go-chain/internal/users"
)

type Blockchain struct {
	Blocks []*Block
	Users  map[string]*users.User
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

func NewBlockchain() *Blockchain {
	genesisBlock := NewGenesisBlock()
	return &Blockchain{[]*Block{genesisBlock}, make(map[string]*users.User)}
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

// FindUnspentTransactions finds transactions with unspent outputs
func (bc *Blockchain) FindUnspentTransactions(pubKeyHash []byte) []transactions.Transaction {
	var unspentTXs []transactions.Transaction
	spentTXOs := make(map[string][]int)

	fmt.Println("Debug: Starting FindUnspentTransactions")

	for _, block := range bc.Blocks {
		fmt.Printf("Debug: Checking block with hash %x\n", block.Hash)
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
			fmt.Printf("Debug: Checking transaction %x\n", tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							fmt.Printf("Debug: Output %d of transaction %x is spent\n", outIdx, tx.ID)
							continue Outputs
						}
					}
				}

				if out.IsLockedWithKey(pubKeyHash) {
					unspentTXs = append(unspentTXs, *tx)
					fmt.Printf("Debug: Output %d of transaction %x is unspent and locked with key\n", outIdx, tx.ID)
				}
			}

			if !tx.IsCoinbase() {
				for _, in := range tx.Vin {
					if in.UsesKey(pubKeyHash) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
						fmt.Printf("Debug: Input of transaction %x uses key and spends output %d of transaction %s\n", tx.ID, in.Vout, inTxID)
					}
				}
			}
		}
	}

	fmt.Println("Debug: Completed FindUnspentTransactions")
	return unspentTXs
}

// FindUTXO returns unspent transaction outputs for a given public key hash
func (bc *Blockchain) FindUTXO(pubKeyHash []byte) []transactions.TransactionOutput {
	var UTXOs []transactions.TransactionOutput
	unspentTransactions := bc.FindUnspentTransactions(pubKeyHash)

	fmt.Printf("Debug: Unspent transactions: %v\n", unspentTransactions)

	for _, tx := range unspentTransactions {
		fmt.Printf("Debug: Checking unspent transaction %x\n", tx.ID)
		for outIdx, out := range tx.Vout {
			if out.IsLockedWithKey(pubKeyHash) {
				fmt.Printf("Debug: Found UTXO in transaction %x at output %d\n", tx.ID, outIdx)
				outCopy := out
				outCopy.Txid = tx.ID
				outCopy.Vout = outIdx
				UTXOs = append(UTXOs, outCopy)
			}
		}
	}

	fmt.Printf("Debug: Total UTXOs found: %d\n", len(UTXOs))
	return UTXOs
}

// GetTransactionsByPubKeyHash retrieves all transactions involving the given public key hash.
func (bc *Blockchain) GetTransactionsByPubKeyHash(pubKeyHash []byte) []*transactions.Transaction {
	var relevantTransactions []*transactions.Transaction

	fmt.Printf("Debug: Starting GetTransactionsByPubKeyHash for pubKeyHash: %x\n", pubKeyHash)

	for _, block := range bc.Blocks {
		fmt.Printf("Debug: Checking block with hash %x\n", block.Hash)
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
			involvesPubKeyHash := false

			// Check if the public key hash is involved in the outputs
			for _, out := range tx.Vout {
				if out.IsLockedWithKey(pubKeyHash) {
					involvesPubKeyHash = true
					break
				}
			}

			// Check if the public key hash is involved in the inputs
			for _, in := range tx.Vin {
				if in.UsesKey(pubKeyHash) {
					involvesPubKeyHash = true
					break
				}
			}

			if involvesPubKeyHash {
				relevantTransactions = append(relevantTransactions, tx)
				fmt.Printf("Debug: Transaction %s involves pubKeyHash %x\n", txID, pubKeyHash)
			}
		}
	}

	fmt.Printf("Debug: Found %d relevant transactions for pubKeyHash: %x\n", len(relevantTransactions), pubKeyHash)
	return relevantTransactions
}
