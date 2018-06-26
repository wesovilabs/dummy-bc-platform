package api

import (
	"net/http"
	"encoding/json"
	"io"
	"dummy-bc-platform/database"
)

func getBC(w http.ResponseWriter, r *http.Request) {
	blockchain := database.GetBlockchain()
	bytes, err := json.MarshalIndent(blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
	w.WriteHeader(http.StatusOK)
}
