package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

// This Defines the Block Structure in the Blockchain
type Block struct {
	Index       int
	Timestamp   string
	PayloadData string
	Hash        string
	PrevHash    string
}

// Declare Blockchain using the Block list
var Blockchain []Block

// SHA-256 Hashing
func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + block.PayloadData + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// Create a new Block (Using PrevHash)
func generateBlock(oldBlock Block, PayloadData string) (Block, error) {
	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.PayloadData = PayloadData
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}

// Make sure the Block Index is correct
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// Declare a Channel for handling incoming concurrent Blocks
var bcServer chan []Block

func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	bcServer = make(chan []Block)

	// Create genesis Block
	t := time.Now()
	genesisBlock := Block{0, t.String(), "INIT: This is Genesis Block.", "", ""}
	spew.Dump(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)

	// Run & Serve TCP Server
	server, err := net.Listen("tcp", ":"+os.Getenv("ADDR"))
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	// Create a New Connection when the Server receives a Connection Request
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

	io.WriteString(conn, "Type the Data >> ")

	scanner := bufio.NewScanner(conn)

	go func() {
		for scanner.Scan() {
			payload_data := scanner.Text()

			newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], payload_data)
			if err != nil {
				log.Println(err)
				continue
			}
			if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
				newBlockchain := append(Blockchain, newBlock)
				replaceChain(newBlockchain)
			}

			bcServer <- Blockchain
			io.WriteString(conn, "\nType the Data >> ")
		}
	}()

	go func() {
		for {
			time.Sleep(30 * time.Second)
			output, err := json.Marshal(Blockchain)

			if err != nil {
				log.Fatal(err)
			}
			io.WriteString(conn, string(output))
		}
	}()

	for _ = range bcServer {
		spew.Dump(Blockchain)
	}
}
