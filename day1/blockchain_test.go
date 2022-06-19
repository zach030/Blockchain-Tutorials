package main

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"testing"
)

func TestBlockchain_AddBlock(t *testing.T) {
	blockchain := NewBlockchain()
	blockchain.AddBlock("Transfer 10$ from Alice To Bob")
	blockchain.AddBlock("Transfer 3$ from Bob To Cat")
	blockchain.AddBlock("Transfer 24$ from James To Alice")
	for _, block := range blockchain.blocks {
		fmt.Printf("%s\n",block.Data)
	}

	data1 := []byte("I like donuts")
	data2 := []byte("I like donutsca07ca")
	targetBits := 24
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	fmt.Printf("%x\n", sha256.Sum256(data1))
	fmt.Printf("%64x\n", target)
	fmt.Printf("%x\n", sha256.Sum256(data2))
}
