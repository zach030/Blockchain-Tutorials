package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

const subsidy = 10

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

type TXOutput struct {
	Value        int
	ScriptPubKey string
}

func NewCoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{subsidy, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()

	return &tx
}

// IsCoinbase checks whether the transaction is coinbase
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// SetID sets ID of a transaction
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	accumulated, validOutputs := bc.FindSpendableOutputs(from, amount)
	if accumulated < amount {
		log.Panic("ERROR: Not enough funds")
	}
	for txid, idxs := range validOutputs {
		txID, _ := hex.DecodeString(txid)
		for _, idx := range idxs {
			inputs = append(inputs, TXInput{
				Txid:      txID,
				Vout:      idx,
				ScriptSig: from,
			})
		}
	}
	outputs = append(outputs, TXOutput{
		Value:        amount,
		ScriptPubKey: to,
	})
	if accumulated > amount {
		outputs = append(outputs, TXOutput{
			Value:        accumulated - amount,
			ScriptPubKey: from,
		})
	}
	tx := Transaction{Vin: inputs, Vout: outputs}
	tx.SetID()
	return &tx
}
