package main

import (
	"fmt"
	"github.com/NicholasRodrigues/go-chain/internal/blockchain"
)

func main() {
	bc := blockchain.NewBlockchain()

	bc.AddBlock("First Block after Genesis")
	bc.AddBlock("Second Block after Genesis")

	for _, block := range bc.Blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}

	fmt.Println("Is blockchain valid?", bc.IsValid())
}
