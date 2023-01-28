package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	nodeURL = "http://localhost:8545"
	pollInterval = 60 * time.Second
	logFile = "geth-performance.log"
)

type ethBlockNumber struct {
	Jsonrpc string `json:"jsonrpc"`
	Id int `json:"id"`
	Result string `json:"result"`
}

type ethGasPrice struct {
	Jsonrpc string `json:"jsonrpc"`
	Id int `json:"id"`
	Result string `json:"result"`
}

type netPeerCount struct {
	Jsonrpc string `json:"jsonrpc"`
	Id int `json:"id"`
	Result string `json:"result"`
}

func main() {
	for {
		// Get the current block number
		blockNumber := getBlockNumber()

		// Get the current gas price
		gasPrice := getGasPrice()

		// Get the current node peer count
		peerCount := getPeerCount()

		// Log the performance data
		logPerformance(blockNumber, gasPrice, peerCount)

		// Sleep for the specified interval
		time.Sleep(pollInterval)
	}
}

func getBlockNumber() string {
	blockNumber := &ethBlockNumber{}
	body := getJSONRPC("eth_blockNumber", []interface{}{})
	json.Unmarshal(body, &blockNumber)
	return blockNumber.Result
}

func getGasPrice() string {
	gasPrice := &ethGasPrice{}
	body := getJSONRPC("eth_gasPrice", []interface{}{})
	json.Unmarshal(body, &gasPrice)
	return gasPrice.Result
}

func getPeerCount() string {
	peerCount := &netPeerCount{}
	body := getJSONRPC("net_peerCount", []interface{}{})
	json.Unmarshal(body, &peerCount)
	return peerCount.Result
}

func getJSONRPC(method string, params []interface{}) []byte {
	// Create the JSON-RPC request
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"method": method,
		"params": params,
		"id": 1,
	}

	// Send the request
	response, _ := http.Post(nodeURL, "application/json", json.NewEncoder(request))
	body, _ := ioutil.ReadAll(response.Body)
	return body
}

func logPerformance(blockNumber, gasPrice, peerCount string) {
	// Log the performance data
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	logString := fmt.Sprintf("Block number: %s, Gas price: %s, Peer count: %s, Timestamp: %s", blockNumber, gasPrice, peerCount, time.Now().String())
	if _, err := f.WriteString(logString + "\n"); err != nil {
		log.Println(err)
	}
}
