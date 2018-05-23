package main

/*
Packages imported but not used will be deleted when you save the file if go.formatTool in user settings   is configured to use goreturns. Go to File\Preferences\Settings and put the following in the workspace settings to overwrite:
{
    "go.formatTool": "gofmt"
}
You will also see a file named settings.json created under .vscode folder if this is your very first time to overwrite a setting.

*/

import (
/*
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
*/
)

//Block ...
type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
}

func main() {

}
