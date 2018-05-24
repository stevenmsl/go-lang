package main

/*
Packages imported but not used will be deleted when you save the file if go.formatTool in user settings   is configured to use goreturns. Go to File\Preferences\Settings and put the following in the workspace settings to overwrite:
{
    "go.formatTool": "gofmt"
}
You will also see a file named settings.json created under .vscode folder if this is your very first time to overwrite a setting.

*/

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
	/*	"encoding/json"
		"io"
		"log"
		"net/http"
		"os"


		"github.com/davecgh/go-spew/spew"
		"github.com/gorilla/mux"
		"github.com/joho/godotenv"
	*/)

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

func main() {

}
