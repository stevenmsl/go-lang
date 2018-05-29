package main

/*
The .env file needs to be in the package folder (src\blockchain) so it can be located when you start debugging.
*/

/*
Set the breakpoints before you start debugging.
*/

/*
The .env file needs to be in the package folder (src\blockchain) so it can be located when you start debugging.
Packages imported but not used will be deleted when you save the file if go.formatTool in user settings   is configured to use goreturns. Go to File\Preferences\Settings and put the following in the workspace settings to overwrite:
{
    "go.formatTool": "gofmt"
}
You will also see a file named settings.json created under .vscode folder if this is your very first time to overwrite a setting.

*/

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

//Block ...
type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
}

//Blockchain ...
var Blockchain []Block

//Message ...
type Message struct {
	BPM int
}

func generateBlock(oldBlock Block, BPM int) (Block, error) {
	var newBlock Block
	t := time.Now()
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)
	return newBlock, nil
}

func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hased := h.Sum(nil)
	return hex.EncodeToString(hased)
}

func isBlockValid(newBlock Block, oldBlock Block) bool {
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

func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

func run() error {
	mux := makeMuxRouter()
	httpAddr := os.Getenv("ADDR")
	log.Println("Listening on ", os.Getenv("ADDR"))
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 2^20
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockChain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

func handleGetBlockChain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var m Message
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, r.Body)
		return
	}
	//A defer statement pushes a function call onto a list.
	//The list of saved calls is executed after the surrounding function returns.
	//What this means is that r.Body.Close will be executed at the very last after other statements have been executed.
	defer r.Body.Close()

	newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], m.BPM)
	if err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, m)
		return
	}
	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		newBlockchain := append(Blockchain, newBlock)
		replaceChain(newBlockchain)
		spew.Dump(Blockchain)
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)

}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int,
	payload interface{} /* this means any type as all types implement the empty interface */) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		t := time.Now()
		genesisBlock := Block{0, t.String(), 0, "", ""}
		genesisBlock.Hash = calculateHash(genesisBlock)
		spew.Dump(genesisBlock)
		Blockchain = append(Blockchain, genesisBlock)
	}()

	log.Fatal(run())

}
