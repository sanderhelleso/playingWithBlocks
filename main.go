package playingWithBlocks

import (
	"crypto/sha256"
	"encoding/hex"
	time2 "time"
	"os"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"encoding/json"
	"io"
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

// do block validation
func isBlockValid(newBlock, oldBlock Block) bool {
	// check incrementing
	if oldBlock.Index + 1 != newBlock.Index {
		return false
	}

	// check previous hash
	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	// check current hash
	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// replace the current chain
func replaceChain(newBlock []Block) {
	if len(newBlock) > len(Blockchain) {
		Blockchain = newBlock
	}
}

// HTTP server
func run() error {
	// set router
	mux := makeMuxRouter()
	// get port from .env
	httpAddr := os.Getenv("ADDR")
	log.Println("Listening on ", os.Getenv(httpAddr))

	// server setup
	server := &http.Server{
		Addr: ":" + httpAddr,
		Handler: mux,
		ReadTimeout: 10 * time2.Second,
		WriteTimeout: 10 * time2.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// check for error
	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// create router
func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	// GET request
	muxRouter.HandleFunc("/", handleGetBlockChain).Methods("GET")
	// Post request
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return  muxRouter
}

// GET handler
func handleGetBlockChain(w http.ResponseWriter, r * http.Request)  {
	// write back full blockchain in JSON format
	bytes, err := json.MarshalIndent(Blockchain, "", " ")

	// check for server error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write to port
	io.WriteString(w, string(bytes))
}

// value function
type Value struct {
	Value int
}




