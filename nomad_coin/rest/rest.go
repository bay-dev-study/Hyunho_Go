package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"nomad_coin/blockchain"
	"nomad_coin/p2p"
	"nomad_coin/utils"
	"nomad_coin/wallet"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var rootUrlWithPort string
var PortInString string

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

type addPeerPayload struct {
	Address string `json:"address"`
	Port    string `json:"port"`
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
	blockchain.WriteBlockchainToJsonEncoder(jsonEncoder)
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
	blockchain.GetBlockchain().CreateNewBlockFromTx()
	p2p.BroadcastNewBlock(blockchain.GetNewestBlock())
	rw.WriteHeader(http.StatusCreated)
}

func handleMempool(rw http.ResponseWriter, r *http.Request) {
	utils.ErrHandler(json.NewEncoder(rw).Encode(blockchain.GetMempoolTx()))
}

func handleTransactions(rw http.ResponseWriter, r *http.Request) {
	var payload addTxPayload
	utils.ErrHandler(json.NewDecoder(r.Body).Decode(&payload))
	tx, err := blockchain.MakeTx(wallet.GetWallet().Address, payload.To, payload.Amount)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(errorMessage{err.Error()})
		return
	}
	blockchain.GetMempool().AddTxToMempool(tx)
	p2p.BroadcastNewTransaction(tx)
	rw.WriteHeader(http.StatusCreated)
}

func handlerWebsocketUpgrade(rw http.ResponseWriter, r *http.Request) {
	openPort := r.URL.Query().Get("openPort")
	ip := utils.Splitter(r.RemoteAddr, ":", 0)

	var upgrader = websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return openPort != "" && ip != ""
	}
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.ErrHandler(err)
	p2p.InitPeer(conn, ip, openPort)
}

func handlePeer(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var payload addPeerPayload
		json.NewDecoder(r.Body).Decode(&payload)
		p2p.AddPeer(payload.Address, payload.Port, PortInString[1:], true)
		rw.WriteHeader(http.StatusOK)
	case "GET":
		p2p.WritePeersToJsonEncoder(json.NewEncoder(rw))
		rw.WriteHeader(http.StatusOK)
	}
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
	router.HandleFunc("/ws", handlerWebsocketUpgrade).Methods("GET")
	router.HandleFunc("/peer", handlePeer).Methods("GET", "POST")

	PortInString = fmt.Sprintf(":%d", port)
	rootUrlWithPort = fmt.Sprintf("http://localhost%s", PortInString)
	blockchain.SetBlockchainDatabaseFileName(strconv.Itoa(port))

	fmt.Printf("Rest server listening on %s\n", rootUrlWithPort)
	log.Fatal(http.ListenAndServe(PortInString, router))
}
