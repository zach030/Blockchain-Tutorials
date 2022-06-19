package main

import (
	"fmt"
	"testing"
)

func TestNewBlockchain(t *testing.T) {
	blockchain := NewBlockchain()
	blockchain.AddBlock("Transfer 10$ from Alice To Bob")
	blockchain.AddBlock("Transfer 3$ from Bob To Cat")
	blockchain.AddBlock("Transfer 24$ from James To Alice")
	iter := blockchain.Iterator()
	var block *Block
	for i := 0; i < 3; i++ {
		block = iter.Next()
		fmt.Printf("Prev Hash:%x\nData: %s\nHash:%x\n\n", block.PrevBlockHash, block.Data, block.Hash)
	}
}
