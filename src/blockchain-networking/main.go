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
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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
var bcServer chan []Block
var mutex = &sync.Mutex{}

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

//take the longest chain
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
	t := time.Now()
	bcServer = make(chan []Block)
	genesisBlock := Block{0, t.String(), 0, "", ""}
	genesisBlock.Hash = calculateHash(genesisBlock)
	spew.Dump(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)

	server, err := net.Listen("tcp", ":"+os.Getenv("TCP_ADDR"))

	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()
	for { //Use infinite loop to accept the new connections. Use go keyword (async call) so it won’t clog up the for loop.
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
	/*
		Why the for scanner.Scan() loop needs to be tucked away in its own Go routine so
		it can run concurrently and separately from other connections?
		Do we not call the handle Conn using Go routine already?
		Answer: I think even inside the same connection you are entering new BPMs and
		receiving broadcast at the same time both using infinite loop.
		Make sense to use Go routine in both cases not to clog up each other.
	*/
	go func() { //Asking client to enter a new BPM
		for scanner.Scan() {
			bpm, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Printf("%v not a number: %v", scanner.Text(), err)
				continue
			}
			newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], bpm)
			if err != nil {
				log.Println(err)
				continue
			}
			if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
				newBlockchain := append(Blockchain, newBlock)
				replaceChain(newBlockchain)
			}
			bcServer <- Blockchain
			io.WriteString(conn, "\n Enter a new BPM:") //prompt user to enter new BPM
		}
	}()

	go func() { //Broadcast the blockchain back to the client
		for {
			time.Sleep(30 * time.Second)
			mutex.Lock()
			output, err := json.Marshal(Blockchain)
			if err != nil {
				log.Fatal(err)
			}
			mutex.Unlock()
			io.WriteString(conn, string(output))

		}
	}()

	//range over channels
	//acting as a receiver
	for _ = range bcServer {
		spew.Dump(Blockchain)
	}

}
