package blocks

import (
	"time"
)

type Block struct {
	// Index position of the data record in the blockchain
	Index int `json:"index"`
	// Timestamp automatically determined and is the time the data is written
	Timestamp time.Time `json:"timestamp"`
	// BPM beats per minute, is your pulse rate
	BPM int `json:"bpm"`
	// Hash SHA256 identifier representing this data record
	Hash string `json:"hash"`
	// PreviousHash  SHA256 identifier of the previous record in the chain
	PreviousHash string `json:"previous_hash"`
}

// BlockChain Array of blocks
type BC []*Block


