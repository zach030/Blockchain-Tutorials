package main

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct {
	TimeStamp int64
	Data []byte
	PrevBlockHash []byte
	Hash []byte
}

func NewBlock(data string, prevBlockHash []byte)*Block{
	block := &Block{
		TimeStamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Hash:          make([]byte,0),
	}
	block.SetHash()
	return block
}

func NewGenesisBlock()*Block{
	b := &Block{
		TimeStamp:     time.Now().Unix(),
		Data:          []byte("Hello BlockChain"),
		PrevBlockHash: nil,
		Hash:          nil,
	}
	b.SetHash()
	return b
}

// hash = sha256(timestamp, data, prevHash)
func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.TimeStamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}