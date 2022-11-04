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

type errorMessage struct {
	ErrorMessage string `json:"errorMessage"`
}

type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type addTxPayload struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

func getDocument() []*documentData {
	documentDataAll := make([]*documentData, 0)
	documentDataAll = append(documentDataAll, &documentData{"/", "GET", "API Document", ""})
	documentDataAll = append(documentDataAll, &documentData{"/status", "GET", "Get Blockchain Status", ""})
	documentDataAll = append(documentDataAll, &documentData{"/blocks", "GET", "Get All Blocks Info", ""})
	documentDataAll = append(documentDataAll, &documentData{"/blocks", "POST", "Add A Block", "{data: data_to_add_in_block}"})
	documentDataAll = append(documentDataAll, &documentData{"/blocks/{hash}", "GET", "See A Block", ""})
	return documentDataAll
}

func getAllBlocks() []*blockchain.Block {
	return blockchain.AllBlocks()
}

func handleRoot(rw http.ResponseWriter, r *http.Request) {
	document := getDocument()
	rw.Header().Add("content-type", "application/json")
	utils.ErrHandler(json.NewEncoder(rw).Encode(&document))
}

func handleStatus(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("content-type", "application/json")
	jsonEncoder := json.NewEncoder(rw)
	jsonEncoder.Encode(blockchain.GetBlockchain())
}

func handleBlocks(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("content-type", "application/json")
	jsonEncoder := json.NewEncoder(rw)
	jsonEncodedAllBlocksInfo := getAllBlocks()
	jsonEncoder.Encode(jsonEncodedAllBlocksInfo)
}

func handleBlocksByHash(rw http.ResponseWriter, r *http.Request) {
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

func handleBalance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	total := r.URL.Query().Get("total")

	switch total {
	case "true":
		amount := blockchain.BalanceByAddress(address)
		json.NewEncoder(rw).Encode(balanceResponse{address, amount})
	default:
		utils.ErrHandler(json.NewEncoder(rw).Encode(blockchain.GetUTxOfAddress(address)))
	}
}

func handleConfirm(rw http.ResponseWriter, r *http.Request) {
	blockchain.GetBlockchain().ConfirmBlock()
	rw.WriteHeader(http.StatusCreated)
}

func handleMempool(rw http.ResponseWriter, r *http.Request) {
	utils.ErrHandler(json.NewEncoder(rw).Encode(blockchain.GetMempool().Txs))
}

func handleTransactions(rw http.ResponseWriter, r *http.Request) {
	var payload addTxPayload
	utils.ErrHandler(json.NewDecoder(r.Body).Decode(&payload))
	err := blockchain.GetMempool().AddTx(payload.From, payload.To, payload.Amount)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(errorMessage{err.Error()})
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

func Start(port int) {
	router := mux.NewRouter()
	router.HandleFunc("/", handleRoot).Methods("GET")
	router.HandleFunc("/status", handleStatus).Methods("GET")
	router.HandleFunc("/blocks", handleBlocks).Methods("GET")
	router.HandleFunc("/confirm", handleConfirm).Methods("GET")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", handleBlocksByHash).Methods("GET")
	router.HandleFunc("/balance/{address}", handleBalance)
	router.HandleFunc("/mempool", handleMempool)
	router.HandleFunc("/transactions", handleTransactions).Methods("POST")

	portInString := fmt.Sprintf(":%d", port)
	rootUrlWithPort = fmt.Sprintf("http://localhost%s", portInString)

	fmt.Printf("Rest server listening on %s\n", rootUrlWithPort)
	log.Fatal(http.ListenAndServe(portInString, router))
}
