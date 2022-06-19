package main

type Blockchain struct {
	blocks []*Block
}

func NewBlockchain()*Blockchain{
	bc := &Blockchain{}
	bc.blocks = make([]*Block,0)
	bc.blocks = append(bc.blocks, NewGenesisBlock())
	return bc
}

func (b *Blockchain) AddBlock(data string) {
	prevBlock := b.blocks[len(b.blocks)-1]
	block := NewBlock(data,prevBlock.PrevBlockHash)
	b.blocks = append(b.blocks, block)
}
