package main

import (
	"bufio"
	"dummy-bc-platform/blocks"
	"dummy-bc-platform/database"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

const serverAddressKey = "SERVER_ADDRESS"

// bcServer handles incoming concurrent Blocks
var bcServer chan *blocks.BC
var mutex = &sync.Mutex{}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	genesisBlock := &blocks.Block{0, time.Now(), 0, "", ""}
	spew.Dump(genesisBlock)
	database.AddBlock(genesisBlock)
	startServer()
}

func startServer() error {
	bcServer = make(chan *blocks.BC)
	address := os.Getenv(serverAddressKey)
	log.Println("Listening on ", address)
	server, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	io.WriteString(conn, "Enter a new BPM:")

	scanner := bufio.NewScanner(conn)

	go func() {
		for scanner.Scan() {
			bpm, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Printf("%v not a number: %v", scanner.Text(), err)
				continue
			}
			lastBlock := database.GetLastBlock()
			block, err := lastBlock.NewBlock(bpm)
			if err != nil {
				log.Println(err)
				continue
			}
			blockchain := database.GetBlockchain()
			if block.IsValid(lastBlock) {
				fmt.Println("Block is valid.")
				newBlockchain := append(*blockchain, block)
				database.UpdateBlockchain(&newBlockchain)
				spew.Dump(blockchain)
			}
			bcServer <- blockchain
			io.WriteString(conn, "\nEnter a new BPM:")
		}
	}()

	// simulate receiving broadcast
	go func() {

		for {
			time.Sleep(30 * time.Second)
			mutex.Lock()
			blockchain := database.GetBlockchain()
			output, err := json.Marshal(blockchain)
			if err != nil {
				log.Fatal(err)
			}
			mutex.Unlock()
			io.WriteString(conn, string(output))
		}
	}()

	for range bcServer {
		blockchain := database.GetBlockchain()
		spew.Dump(blockchain)
	}

}
