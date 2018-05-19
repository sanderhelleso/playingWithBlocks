package playingWithBlocks

import (
	"crypto/sha256"
	"encoding/hex"
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


