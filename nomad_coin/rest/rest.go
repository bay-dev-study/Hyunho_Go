package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"nomad_coin/blockchain"
	"nomad_coin/utils"

	"github.com/gorilla/mux"
)

var rootUrlWithPort string

type url string

type documentData struct {
	Url         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("%s%s", rootUrlWithPort, u)
	return []byte(url), nil
}

type blockPayload struct {
	Data string `json:"data"`
}

type errorMessage struct {
	ErrorMessage string `json:"errorMessage"`
}

func getDocument() []*documentData {
	documentDataAll := make([]*documentData, 0)
	documentDataAll = append(documentDataAll, &documentData{"/", "GET", "API Document", ""})
	documentDataAll = append(documentDataAll, &documentData{"/blocks", "GET", "Get All Blocks Info", ""})
	documentDataAll = append(documentDataAll, &documentData{"/blocks", "POST", "Add A Block", "{data: data_to_add_in_block}"})
	documentDataAll = append(documentDataAll, &documentData{"/blocks/{hash}", "GET", "See A Block", ""})
	return documentDataAll
}

func getAllBlocks() []*blockchain.Block {
	return blockchain.GetBlockchain().AllBlocks()
}

func handleRoot(rw http.ResponseWriter, r *http.Request) {
	document := getDocument()
	rw.Header().Add("content-type", "application/json")
	utils.ErrHandler(json.NewEncoder(rw).Encode(&document))
}

func handleBlocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rw.Header().Add("content-type", "application/json")
		jsonEncoder := json.NewEncoder(rw)
		jsonEncodedAllBlocksInfo := getAllBlocks()
		jsonEncoder.Encode(jsonEncodedAllBlocksInfo)
	case "POST":
		var blockDataToAdd blockPayload
		rw.Header().Add("content-type", "application/json")
		utils.ErrHandler(json.NewDecoder(r.Body).Decode(&blockDataToAdd))
		blockchain.GetBlockchain().CreateBlockAndSave(blockDataToAdd.Data)
		rw.WriteHeader(http.StatusCreated)
	}
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	blockHashInRawUrl := mux.Vars(r)
	blockHash := blockHashInRawUrl["hash"]

	blockDataMatchesHash, err := blockchain.GetBlockByHash(blockHash)
	if err != nil {
		errorMessage := errorMessage{"block not found"}
		json.NewEncoder(rw).Encode(&errorMessage)
	} else {
		json.NewEncoder(rw).Encode(&blockDataMatchesHash)
	}
}

func Start(port int) {
	router := mux.NewRouter()
	router.HandleFunc("/", handleRoot).Methods("GET")
	router.HandleFunc("/blocks", handleBlocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", blocks).Methods("GET")

	portInString := fmt.Sprintf(":%d", port)
	rootUrlWithPort = fmt.Sprintf("http://localhost%s", portInString)

	fmt.Printf("Rest server listening on %s\n", rootUrlWithPort)
	log.Fatal(http.ListenAndServe(portInString, router))
}
