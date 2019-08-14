package main

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
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

func main() {

}
