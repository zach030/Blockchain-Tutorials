package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type CLI struct {
	bc *Blockchain
}

func (c *CLI) Run() {
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	addBlockData := addBlockCmd.String("data", "", "block data")
	var err error
	switch os.Args[1] {
	case "addblock":
		err = addBlockCmd.Parse(os.Args[2:])
	case "printchain":
		err = printChainCmd.Parse(os.Args[2:])
	default:
		os.Exit(1)
	}
	if err != nil {
		panic(err)
	}
	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		c.addBlock(*addBlockData)
	}
	if printChainCmd.Parsed() {
		c.printChain()
	}
}

func (c *CLI) addBlock(data string) {
	c.bc.AddBlock(data)
	fmt.Println("Success!")
}

func (c *CLI) printChain() {
	bci := c.bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := newProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
