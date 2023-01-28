package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const nodeURL = "https://mainnet.infura.io/v3/6ebc20d69d504a9c961be21a55cfe2e9"
const logFile = "node_performance.log"

type request struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      string        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

func main() {
	for {
		// Get the current block number
		requestData := request{JSONRPC: "2.0", ID: "1", Method: "eth_blockNumber", Params: []interface{}{}}
		requestBody := new(bytes.Buffer)
		json.NewEncoder(requestBody).Encode(requestData)
		response, _ := http.Post(nodeURL, "application/json", requestBody)

		var result map[string]interface{}
		json.NewDecoder(response.Body).Decode(&result)
		blockNumber := result["result"].(string)

		// Get the current gas price
		requestData = request{JSONRPC: "2.0", ID: "1", Method: "eth_gasPrice", Params: []interface{}{}}
		requestBody = new(bytes.Buffer)
		json.NewEncoder(requestBody).Encode(requestData)
		response, _ = http.Post(nodeURL, "application/json", requestBody)

		json.NewDecoder(response.Body).Decode(&result)
		gasPrice := result["result"].(string)

		// Get the current peer count
		requestData = request{JSONRPC: "2.0", ID: "1", Method: "net_peerCount", Params: []interface{}{}}
		requestBody = new(bytes.Buffer)
		json.NewEncoder(requestBody).Encode(requestData)
		response, _ = http.Post(nodeURL, "application/json", requestBody)

		json.NewDecoder(response.Body).Decode(&result)
		peerCount := result["result"].(string)

		logPerformance(blockNumber, gasPrice, peerCount)

		time.Sleep(60 * time.Second)
	}
}

func logPerformance(blockNumber, gasPrice, peerCount string) {
    // Log the performance data
    f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Println(err)
        return
    }
    defer f.Close()

    logString := fmt.Sprintf("Block number: %s, Gas price: %s, Peer count: %s, Timestamp: %s\n", blockNumber, gasPrice, peerCount, time.Now().String())
    if _, err := f.WriteString(logString); err != nil {
        log.Println(err)
    }
}
