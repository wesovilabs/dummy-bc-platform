package api

import (
	"dummy-bc-platform/database"
	"encoding/json"
	"io"
	"net/http"
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
