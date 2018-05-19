package playingWithBlocks

import (
	"crypto/sha256"
	"encoding/hex"
	time2 "time"
)

// blockChain structure
type Block struct {
	Index int
	Timestamp string
	Value int
	Hash string
	PrevHash string
}

// the main blockchain
var Blockchain []Block

// calculate the hash of each block
func calculateHash(block Block) string  {
	record := string(block.Index) + block.Timestamp + string(block.Value) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// generates a new block
func generateBlock(oldBlock Block, Value int) (Block, error) {
	var newBlock Block
	time :=  time2.Now()

	// set values for the block
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = time.String()
	newBlock.Value = Value
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}


