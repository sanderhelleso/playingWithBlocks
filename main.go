package playingWithBlocks

import (
	"crypto/sha256"
	"encoding/hex"
	time2 "time"
	"os"
	"log"
	"net/http"
	"encoding/json"
	"io"

	"github.com/gorilla/mux"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
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

// POST handler
func handleWriteBlock(w http.ResponseWriter, r * http.Request) {
	// value for blockchain
	var m Value

	// check for 400
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}

	defer r.Body.Close()

	// check for 500
	newBlock, err := generateBlock(Blockchain[len(Blockchain) - 1], m.Value)
	if err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, m)
		return
	}

	// check validation of the block
	if isBlockValid(newBlock, Blockchain[len(Blockchain) - 1]) {
		newBlockChain := append(Blockchain, newBlock)
		replaceChain(newBlockChain)
		//pretty prints our structs into the console
		spew.Dump(Blockchain)
	}

	// send the chain in JSON format
	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

// create chain as JSON and respond to request
func respondWithJSON(w http.ResponseWriter, r * http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", " ")
	// check for error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	// write header and response
	w.WriteHeader(code)
	w.Write(response)
}

// main
func main() {
	// allow us to read from .env
	err := godotenv.Load()
	// check for error
	if err != nil {
		log.Fatal(err)
	}

	// init chain
	go func() {
		time := time2.Now()
		genesisBlock := Block{0, time.String(), 0, "", ""}
		spew.Dump(genesisBlock)
		Blockchain = append(Blockchain, genesisBlock)
	}()
	log.Fatal(run())
}



