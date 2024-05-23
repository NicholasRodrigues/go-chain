package blockchain

import (
	"fmt"
	"github.com/NicholasRodrigues/go-chain/internal/transactions"
	"testing"
)

func TestCreateUser(t *testing.T) {
	bc := NewBlockchain()
	username := "Alice"
	password := "password123"

	user := bc.CreateUser(username, password)

	if user.Username != username {
		t.Errorf("Expected username %s, got %s", username, user.Username)
	}
	if user.Password != password {
		t.Errorf("Expected password to be set correctly")
	}
	if user.PrivateKey == nil || user.PublicKey == nil {
		t.Errorf("Expected keys to be generated")
	}
}

func TestGetBalance(t *testing.T) {
	bc := NewBlockchain()
	username := "Alice"
	password := "password123"
	bc.CreateUser(username, password)

	balance, err := bc.GetBalance(username, password)
	if err != nil {
		t.Errorf("Error getting balance: %v", err)
	}
	if balance != 0 {
		t.Errorf("Expected balance 0, got %d", balance)
	}

	err = bc.AddFunds(username, password, 100)
	if err != nil {
		t.Errorf("Error adding funds: %v", err)
	}

	balance, err = bc.GetBalance(username, password)
	if err != nil {
		t.Errorf("Error getting balance: %v", err)
	}
	if balance != 100 {
		t.Errorf("Expected balance 100, got %d", balance)
	}
}

func TestAddFunds(t *testing.T) {
	bc := NewBlockchain()
	username := "Alice"
	password := "password123"
	bc.CreateUser(username, password)

	err := bc.AddFunds(username, password, 50)
	if err != nil {
		t.Errorf("Error adding funds: %v", err)
	}

	balance, err := bc.GetBalance(username, password)
	if err != nil {
		t.Errorf("Error getting balance: %v", err)
	}
	if balance != 50 {
		t.Errorf("Expected balance 50, got %d", balance)
	}
}

func TestSendMoney(t *testing.T) {
	bc := NewBlockchain()
	fromUsername := "Alice"
	fromPassword := "password123"
	toUsername := "Bob"
	toPassword := "password456"
	bc.CreateUser(fromUsername, fromPassword)
	bc.CreateUser(toUsername, toPassword)

	// Add initial funds to Alice's account
	err := bc.AddFunds(fromUsername, fromPassword, 100)
	if err != nil {
		t.Errorf("Error adding funds: %v", err)
	}

	// Send money from Alice to Bob
	err = bc.SendMoney(fromUsername, fromPassword, toUsername, 30)
	if err != nil {
		t.Errorf("Error sending money: %v", err)
	}

	// Check Alice's balance
	balance, err := bc.GetBalance(fromUsername, fromPassword)
	if err != nil {
		t.Errorf("Error getting balance: %v", err)
	}
	if balance != 70 {
		t.Errorf("Expected balance 70, got %d", balance)
	}

	// Check Bob's balance
	balance, err = bc.GetBalance(toUsername, toPassword)
	if err != nil {
		t.Errorf("Error getting balance: %v", err)
	}
	if balance != 30 {
		t.Errorf("Expected balance 30, got %d", balance)
	}
}

func TestGetTransactionsByPubKeyHash(t *testing.T) {
	// Create a new blockchain
	bc := NewBlockchain()

	// Create users
	userAlice := bc.CreateUser("Alice", "password123")
	userBob := bc.CreateUser("Bob", "password456")

	// Add funds to Alice
	err := bc.AddFunds("Alice", "password123", 100)
	if err != nil {
		t.Fatalf("Failed to add funds to Alice: %v", err)
	}

	// Send money from Alice to Bob
	err = bc.SendMoney("Alice", "password123", "Bob", 50)
	if err != nil {
		t.Fatalf("Failed to send money from Alice to Bob: %v", err)
	}

	// Get Alice's public key hash
	pubKeyHashAlice := transactions.HashPubKey(userAlice.PublicKey.Bytes())

	// Retrieve and print all transactions involving Alice's public key hash
	txs := bc.GetTransactionsByPubKeyHash(pubKeyHashAlice)
	fmt.Printf("Transactions involving Alice:\n")
	for _, tx := range txs {
		fmt.Println(tx.String())
	}

	// Get Bob's public key hash
	pubKeyHashBob := transactions.HashPubKey(userBob.PublicKey.Bytes())

	// Retrieve and print all transactions involving Bob's public key hash
	txs = bc.GetTransactionsByPubKeyHash(pubKeyHashBob)
	fmt.Printf("Transactions involving Bob:\n")
	for _, tx := range txs {
		fmt.Println(tx.String())
	}
}
