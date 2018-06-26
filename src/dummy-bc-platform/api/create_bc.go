package api

import (
	"net/http"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"dummy-bc-platform/database"
	"dummy-bc-platform/blocks"
)

func createBC(w http.ResponseWriter, r *http.Request) {
	var m Message
	lastBlock := database.GetLastBlock()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()
	var block *blocks.Block
	block, err := lastBlock.NewBlock(m.BPM)
	if err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, m)
		return
	}
	blockchain := database.GetBlockchain()
	if block.IsValid(lastBlock) {
		blocks := append(*blockchain, block)
		blockchain.ReplaceChain(&blocks)
		spew.Dump(blockchain)
	}
	respondWithJSON(w, r, http.StatusCreated, block)
	w.WriteHeader(http.StatusCreated)
}
