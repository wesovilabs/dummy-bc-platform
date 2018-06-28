package database

import (
	"dummy-bc-platform/blocks"
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

// blockchain global variable
var blockchain *blocks.BC

// GetLastBlock returns the last block in the chain
func GetLastBlock() *blocks.Block {
	lastBlock := (*blockchain)[len(*blockchain)-1]
	return lastBlock
}

// GetBlockchain return the chain
func GetBlockchain() *blocks.BC {
	return blockchain
}

// AddBlock adds a new block into the chain
func AddBlock(block *blocks.Block) {
	if blockchain != nil {
		fmt.Println("Bc is not nil")
		newBlockchain := append(*blockchain, block)
		blockchain = &newBlockchain
	} else {
		fmt.Println("Bc is  nil")
		blockchain = &blocks.BC{block}
	}
	spew.Dump(blockchain)
}

// UpdateBlockchain updates the chain
func UpdateBlockchain(bc *blocks.BC) {
	if len(*bc) > len(*blockchain) {
		blockchain = bc
	}

}
