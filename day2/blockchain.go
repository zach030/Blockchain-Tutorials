package main

type Blockchain struct {
	blocks []*Block
}

func NewBlockchain() *Blockchain {
	bc := &Blockchain{}
	bc.blocks = make([]*Block, 0)
	block := NewGenesisBlock()
	bc.blocks = append(bc.blocks, block)
	return bc
}

func (b *Blockchain) AddBlock(data string) {
	prevBlock := b.blocks[len(b.blocks)-1]
	block := NewBlock(data, prevBlock.Hash)
	b.blocks = append(b.blocks, block)
}
