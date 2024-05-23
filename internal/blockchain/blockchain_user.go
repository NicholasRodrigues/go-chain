package blockchain

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/NicholasRodrigues/go-chain/internal/transactions"
	"github.com/NicholasRodrigues/go-chain/internal/users"
)

func (bc *Blockchain) CreateUser(name, password string) *users.User {
	user, err := users.NewUser(name, password)
	if err != nil {
		panic(err)
	}
	bc.Users[name] = user
	return user
}

// GetBalance returns the balance of a user by username and password
func (bc *Blockchain) GetBalance(username, password string) (int, error) {
	user, exists := bc.Users[username]
	if !exists || user.Password != password {
		return 0, errors.New("invalid username or password")
	}

	pubKeyHash := transactions.HashPubKey(user.PublicKey.Bytes())
	fmt.Printf("Debug: Public key hash for user %s: %x\n", username, pubKeyHash)

	utxos := bc.FindUTXO(pubKeyHash)

	balance := 0
	for _, utxo := range utxos {
		fmt.Printf("Debug: Adding UTXO value %d to balance\n", utxo.Value)
		balance += utxo.Value
	}

	fmt.Printf("Debug: Total balance for user %s: %d\n", username, balance)
	return balance, nil
}

// AddFunds adds an arbitrary amount to a user's account for testing purposes
func (bc *Blockchain) AddFunds(username, password string, amount int) error {
	user, exists := bc.Users[username]
	if !exists || user.Password != password {
		return errors.New("invalid username or password")
	}

	bc.Users[username].Balance += amount

	// Use the correct public key hash for the ScriptPubKey
	pubKeyHash := transactions.HashPubKey(user.PublicKey.Bytes())

	// Create a coinbase transaction to add funds
	coinbaseTx := transactions.NewTransaction(
		[]transactions.TransactionInput{{Txid: []byte{}, Vout: -1, ScriptSig: "coinbase"}},
		[]transactions.TransactionOutput{{Value: amount, ScriptPubKey: hex.EncodeToString(pubKeyHash)}},
	)

	// Add the coinbase transaction to a new block
	newBlock := NewBlock([]*transactions.Transaction{coinbaseTx}, bc.Blocks[len(bc.Blocks)-1].Hash)
	bc.Blocks = append(bc.Blocks, newBlock)

	fmt.Println("Coinbase transaction added:", coinbaseTx)
	fmt.Println("New block added with hash:", newBlock.Hash)

	return nil
}

// SendMoney allows a user to send money to another user
func (bc *Blockchain) SendMoney(fromUsername, fromPassword, toUsername string, amount int) error {
	fromUser, fromExists := bc.Users[fromUsername]
	toUser, toExists := bc.Users[toUsername]

	if !fromExists || !toExists || fromUser.Password != fromPassword {
		return errors.New("invalid sender username or password, or recipient username not found")
	}

	fromBalance, _ := bc.GetBalance(fromUsername, fromPassword)
	if fromBalance < amount {
		return errors.New("insufficient funds")
	}

	pubKeyHash := transactions.HashPubKey(fromUser.PublicKey.Bytes())
	utxos := bc.FindUTXO(pubKeyHash)

	var inputs []transactions.TransactionInput
	var outputs []transactions.TransactionOutput

	accumulated := 0
	for _, utxo := range utxos {
		input := transactions.TransactionInput{Txid: utxo.Txid, Vout: utxo.Vout, ScriptSig: ""}
		inputs = append(inputs, input)
		accumulated += utxo.Value
		if accumulated >= amount {
			break
		}
	}

	outputs = append(outputs, transactions.TransactionOutput{Value: amount, ScriptPubKey: hex.EncodeToString(toUser.PublicKey.Bytes())})
	if accumulated > amount {
		outputs = append(outputs, transactions.TransactionOutput{Value: accumulated - amount, ScriptPubKey: hex.EncodeToString(fromUser.PublicKey.Bytes())})
	}

	tx := transactions.NewTransaction(inputs, outputs)
	tx.Sign(fromUser.PrivateKey)

	newBlock := NewBlock([]*transactions.Transaction{tx}, bc.Blocks[len(bc.Blocks)-1].Hash)
	bc.Blocks = append(bc.Blocks, newBlock)

	return nil
}
