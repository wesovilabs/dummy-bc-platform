package blocks

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// NewBlock create new instance of block
func (block *Block) NewBlock(BPM int) (*Block, error) {
	newBlock := &Block{}
	newBlock.Index = block.Index + 1
	newBlock.Timestamp = time.Now()
	newBlock.BPM = BPM
	newBlock.PreviousHash = block.Hash
	newBlock.CalculateHash()
	return newBlock, nil
}

// CalculateHash calculates hash from structure
func (block *Block) CalculateHash() {
	record := string(block.Index) + block.Timestamp.String() + string(block.BPM) + block.PreviousHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	block.Hash = hex.EncodeToString(hashed)
}

// IsValid checks if block is valid
func (block *Block) IsValid(oldBlock *Block) bool {
	if oldBlock.Index+1 != block.Index {
		return false
	}
	if oldBlock.Hash != block.PreviousHash {
		return false
	}
	block.CalculateHash()
	if block.Hash != block.Hash {
		return false
	}
	return true
}

// ReplaceChain replaces the current block chain by the new one
func (bc *BC) ReplaceChain(newBlocks *BC) {
	if len(*newBlocks) > len(*bc) {
		bc = newBlocks
	}
}
