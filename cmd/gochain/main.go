package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"

	"github.com/NicholasRodrigues/go-chain/internal/blockchain"
	"github.com/NicholasRodrigues/go-chain/internal/transactions"
	"github.com/NicholasRodrigues/go-chain/pkg/crypto"
)

func main() {
	bc := blockchain.NewBlockchain()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		printMenu()
		scanner.Scan()
		cmd := scanner.Text()

		switch cmd {
		case "1":
			handleAddBlock(bc, scanner)
		case "2":
			handleViewBlockchain(bc)
		case "3":
			handleValidateBlockchain(bc)
		case "4":
			handleCreateTransaction(scanner)
		case "5":
			handleValidateTransaction(scanner)
		case "6":
			handleFindUTXOs(bc, scanner)
		case "7":
			handleCreateUser(bc, scanner)
		case "8":
			handleGetBalance(bc, scanner)
		case "9":
			handleSendMoney(bc, scanner)
		case "10":
			handleAddFunds(bc, scanner)
		case "11":
			handleExit()
		default:
			fmt.Println("Invalid command. Please try again.")
		}
	}
}

func printMenu() {
	fmt.Println("\nBlockchain CLI")
	fmt.Println("1. Add Block")
	fmt.Println("2. View Blockchain")
	fmt.Println("3. Validate Blockchain")
	fmt.Println("4. Create Transaction")
	fmt.Println("5. Validate Transaction")
	fmt.Println("6. Find UTXOs")
	fmt.Println("7. Create User")
	fmt.Println("8. Get Balance")
	fmt.Println("9. Send Money")
	fmt.Println("10. Add Funds")
	fmt.Println("11. Exit")
	fmt.Print("Enter command: ")
}

func handleAddBlock(bc *blockchain.Blockchain, scanner *bufio.Scanner) {
	fmt.Print("Enter data for the new block: ")
	scanner.Scan()
	data := scanner.Text()

	input := func() string {
		return "input data"
	}
	receive := func() string {
		return data
	}

	blockchain.InputContributionFunction([]byte(data), bc, len(bc.Blocks), input, receive)
	fmt.Println("Block added successfully!")
}

func handleViewBlockchain(bc *blockchain.Blockchain) {
	for i, block := range bc.Blocks {
		fmt.Printf("Block %d: %x\n", i, block.Hash)
		fmt.Printf("Previous Hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Transactions: \n")
		for _, tx := range block.Transactions {
			fmt.Printf("  %s\n", tx.String())
		}
		fmt.Println()
	}
}

func handleValidateBlockchain(bc *blockchain.Blockchain) {
	if blockchain.ChainValidationPredicate(bc) {
		fmt.Println("Blockchain is valid.")
	} else {
		fmt.Println("Blockchain is invalid.")
	}
}

func handleCreateTransaction(scanner *bufio.Scanner) {
	fmt.Print("Enter Txid: ")
	scanner.Scan()
	txid := scanner.Text()

	fmt.Print("Enter Vout: ")
	scanner.Scan()
	vout, _ := strconv.Atoi(scanner.Text())

	fmt.Print("Enter ScriptSig: ")
	scanner.Scan()
	scriptSig := scanner.Text()

	fmt.Print("Enter Value: ")
	scanner.Scan()
	value, _ := strconv.Atoi(scanner.Text())

	fmt.Print("Enter ScriptPubKey: ")
	scanner.Scan()
	scriptPubKey := scanner.Text()

	tx := transactions.NewTransaction(
		[]transactions.TransactionInput{{Txid: []byte(txid), Vout: vout, ScriptSig: scriptSig}},
		[]transactions.TransactionOutput{{Value: value, ScriptPubKey: scriptPubKey}},
	)

	privKey, _ := crypto.NewPrivateKey()
	tx.Sign(privKey)

	fmt.Println("Transaction created successfully!")
	fmt.Println(tx.String())
}

func handleValidateTransaction(scanner *bufio.Scanner) {
	fmt.Print("Enter serialized transaction: ")
	scanner.Scan()
	serializedTx := scanner.Text()
	txBytes, _ := hex.DecodeString(serializedTx)
	tx := transactions.DeserializeTransaction(txBytes)

	utxoSet := make(map[string]transactions.TransactionOutput)
	fmt.Println("Enter UTXO set (enter 'done' to finish):")
	for {
		fmt.Print("Enter UTXO key (txid:vout): ")
		scanner.Scan()
		key := scanner.Text()
		if key == "done" {
			break
		}

		fmt.Print("Enter value: ")
		scanner.Scan()
		value, _ := strconv.Atoi(scanner.Text())

		fmt.Print("Enter ScriptPubKey: ")
		scanner.Scan()
		scriptPubKey := scanner.Text()

		utxoSet[key] = transactions.TransactionOutput{Value: value, ScriptPubKey: scriptPubKey}
	}

	if tx.Validate(utxoSet) {
		fmt.Println("Transaction is valid.")
	} else {
		fmt.Println("Transaction is invalid.")
	}
}

func handleFindUTXOs(bc *blockchain.Blockchain, scanner *bufio.Scanner) {
	fmt.Print("Enter public key hash: ")
	scanner.Scan()
	pubKeyHash, _ := hex.DecodeString(scanner.Text())

	utxos := bc.FindUTXO(pubKeyHash)
	if len(utxos) == 0 {
		fmt.Println("No UTXOs found.")
		return
	}

	fmt.Println("Unspent Transaction Outputs:")
	for _, utxo := range utxos {
		fmt.Printf("Value: %d, ScriptPubKey: %s\n", utxo.Value, utxo.ScriptPubKey)
	}
}

func handleCreateUser(bc *blockchain.Blockchain, scanner *bufio.Scanner) {
	fmt.Print("Enter username: ")
	scanner.Scan()
	username := scanner.Text()

	fmt.Print("Enter password: ")
	scanner.Scan()
	password := scanner.Text()

	user := bc.CreateUser(username, password)
	fmt.Printf("User %s created with public key: %x\n", username, user.PublicKey.Bytes())
}

func handleGetBalance(bc *blockchain.Blockchain, scanner *bufio.Scanner) {
	fmt.Print("Enter username: ")
	scanner.Scan()
	username := scanner.Text()

	fmt.Print("Enter password: ")
	scanner.Scan()
	password := scanner.Text()

	balance, err := bc.GetBalance(username, password)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("User %s has a balance of %d\n", username, balance)
	}
}

func handleSendMoney(bc *blockchain.Blockchain, scanner *bufio.Scanner) {
	fmt.Print("Enter sender's username: ")
	scanner.Scan()
	from := scanner.Text()

	fmt.Print("Enter sender's password: ")
	scanner.Scan()
	fromPassword := scanner.Text()

	fmt.Print("Enter recipient's username: ")
	scanner.Scan()
	to := scanner.Text()

	fmt.Print("Enter amount to send: ")
	scanner.Scan()
	amount, _ := strconv.Atoi(scanner.Text())

	err := bc.SendMoney(from, fromPassword, to, amount)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Transaction successful!")
	}
}

func handleAddFunds(bc *blockchain.Blockchain, scanner *bufio.Scanner) {
	fmt.Print("Enter username: ")
	scanner.Scan()
	username := scanner.Text()

	fmt.Print("Enter password: ")
	scanner.Scan()
	password := scanner.Text()

	fmt.Print("Enter amount to add: ")
	scanner.Scan()
	amount, _ := strconv.Atoi(scanner.Text())

	err := bc.AddFunds(username, password, amount)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Funds added successfully!")
	}
}

func handleExit() {
	fmt.Println("Exiting...")
	os.Exit(0)
}
