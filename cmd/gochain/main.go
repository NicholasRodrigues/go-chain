package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/NicholasRodrigues/go-chain/internal/blockchain"
	"github.com/NicholasRodrigues/go-chain/internal/transactions"
	"github.com/NicholasRodrigues/go-chain/pkg/crypto"
	"os"
	"strconv"
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
	fmt.Println("6. Exit")
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

func handleExit() {
	fmt.Println("Exiting...")
	os.Exit(0)
}
