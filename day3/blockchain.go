package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

type Blockchain struct {
	tip    []byte
	db     *bolt.DB
	blocks []*Block
}

func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil
	}
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		if bucket == nil {
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}
			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}
			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = genesis.Hash
			fmt.Printf("New Blockchain, genesis block hash is:%x\n", tip)
		} else {
			tip = bucket.Get([]byte("l"))
			fmt.Printf("Open Blockchain, current last block hash is:%x\n", tip)
		}
		return nil
	})
	return &Blockchain{tip: tip, db: db}
}

func (b *Blockchain) AddBlock(data string) {
	var lastBlockHash []byte
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		lastBlockHash = bucket.Get([]byte("l"))
		return nil
	})
	block := NewBlock(data, lastBlockHash)
	err = b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		err := bucket.Put(block.Hash, block.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = bucket.Put([]byte("l"), block.Hash)
		if err != nil {
			log.Panic(err)
		}
		b.tip = block.Hash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (b *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{b.tip, b.db}
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block
	err := i.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		encodedBlock := bucket.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	})
	i.currentHash = block.PrevBlockHash
	if err != nil {
		log.Panic(err)
	}
	return block
}
