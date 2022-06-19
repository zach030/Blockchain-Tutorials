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
	for _, block := range blockchain.blocks {
		pow := newProofOfWork(block)
		fmt.Printf("Prev Hash:%x\nData: %s\nHash:%x\nPOW is:%v\n\n", block.PrevBlockHash, block.Data, block.Hash, pow.Validate())
	}
}
