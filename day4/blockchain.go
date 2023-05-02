package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

type Blockchain struct {
	tip    []byte
	db     *bolt.DB
	blocks []*Block
}

func NewBlockchain(address string) *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil
	}
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		if bucket == nil {
			cbtx := NewCoinbaseTx(address, genesisCoinbaseData)
			genesis := NewGenesisBlock(cbtx)
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

func (b *Blockchain) MineBlock(txs []*Transaction) {
	var lastBlockHash []byte
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		lastBlockHash = bucket.Get([]byte("l"))
		return nil
	})
	newBlock := NewBlock(lastBlockHash, txs)
	err = b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		err := bucket.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = bucket.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}
		b.tip = newBlock.Hash
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

func (b *Blockchain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput
	unspentTransactions := b.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

func (b *Blockchain) FindUnspentTransactions(address string) []Transaction {
	unspentTXs := make([]Transaction, 0)
	spentTXOs := make(map[string][]int) // key:txID, value:[]int{idx}
	iter := b.Iterator()
	for {
		block := iter.Next()
		// 遍历区块中的所有交易
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
		Outputs:
			for idx, output := range tx.Vout {
				// 遍历交易的输出
				if spentTXOs[txID] != nil {
					for _, out := range spentTXOs[txID] {
						if out == idx {
							continue Outputs
						}
					}
				}
				if output.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}
			// 排除coinbase tx
			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.Txid)
						// 根据所有的txin记录已使用的输出
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}
		// genesis block
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return unspentTXs
}

func (b *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTx := b.FindUnspentTransactions(address)
	accumulated := 0
loop:
	for _, tx := range unspentTx {
		txID := hex.EncodeToString(tx.ID)
		for idx, output := range tx.Vout {
			if output.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += output.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], idx)
				if accumulated > amount {
					break loop
				}
			}
		}
	}
	return accumulated, unspentOutputs
}
