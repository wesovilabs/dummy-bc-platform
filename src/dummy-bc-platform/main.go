package main

import (
	"net/http"
	"github.com/go-chi/chi"
	"os"
	"log"
	"time"
	"dummy-bc-platform/api"
	"github.com/joho/godotenv"
	"github.com/davecgh/go-spew/spew"
	"dummy-bc-platform/blocks"
	"dummy-bc-platform/database"
)

const serverAddressKey = "SERVER_ADDRESS"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		genesisBlock := &blocks.Block{0, time.Now(), 0, "", ""}
		spew.Dump(genesisBlock)
		database.AddBlock(genesisBlock)
	}()
	r := api.CreateRouter()
	log.Fatal(runServer(r))
}

func runServer(r *chi.Mux) error {

	address := os.Getenv(serverAddressKey)
	log.Println("Listening on ", address)
	server := &http.Server{
		Addr:           address,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return server.ListenAndServe()
}
