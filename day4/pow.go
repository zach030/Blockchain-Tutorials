package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

const targetBits = 24
const maxNonce = math.MaxInt64

type proofOfWork struct {
	block  *Block
	target *big.Int
}

func newProofOfWork(b *Block) *proofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	return &proofOfWork{
		block:  b,
		target: target,
	}
}

func (p *proofOfWork) prepareData(nonce int) []byte {
	datas := bytes.Join([][]byte{p.block.PrevBlockHash,
		p.block.HashTransactions(),
		IntToHex(p.block.TimeStamp),
		IntToHex(targetBits),
		IntToHex(int64(nonce))},
		[]byte{})
	return datas
}

func (p *proofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	fmt.Printf("Mining a new block")
	for nonce < maxNonce {
		data := p.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(p.target) == -1 {
			fmt.Printf("Mint hash is \r%x\n", hash)
			break
		}
		nonce++
	}
	return nonce, hash[:]
}

func (p *proofOfWork) Validate() bool {
	var hashInt big.Int
	data := p.prepareData(p.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	isValid := hashInt.Cmp(p.target) == -1
	return isValid
}

func IntToHex(i int64) []byte {
	return []byte(strconv.FormatInt(i, 10))
}
