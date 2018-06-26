package database

import "dummy-bc-platform/blocks"

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
		newBlockchain := append(*blockchain, block)
		blockchain = &newBlockchain
		return
	}
	blockchain = &blocks.BC{block}
}
