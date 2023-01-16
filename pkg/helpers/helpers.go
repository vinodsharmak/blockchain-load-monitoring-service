package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"path/filepath"
	"runtime"

	"testing"

	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/config"

	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/constants"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/model"
)

func TxPoolContent() (model.Response, error) {
	var response model.Response
	data, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "txpool_content",
		"id":      config.ChainID,
		"params":  []interface{}{},
	})
	if err != nil {
		return response, err
	}

	resp, err := http.Post(config.BlockchainURL, constants.Response_format, strings.NewReader(string(data)))
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func PrepareEmailBodyForTxPoolContent(txpoolContent model.TxPoolContentStuckMail) (string, error) {
	msg := fmt.Sprintf("Pending Count = %v \n", txpoolContent.PendingCount)
	msg = msg + fmt.Sprintf("Queued Count = %v \n{\"Pending Transactions\": [", txpoolContent.QueuedCount)
	for i, txBody := range txpoolContent.PendingContent {
		bytes, err := json.Marshal(txBody)
		if err != nil {
			return msg, err
		}
		if i == len(txpoolContent.PendingContent)-1 {
			msg = msg + string(bytes)
		} else {
			msg = msg + string(bytes) + ",\n"
		}

	}
	msg = msg + "], \n\"Queued Transaction \": ["
	for i, txBody := range txpoolContent.QueuedContent {
		bytes, err := json.Marshal(txBody)
		if err != nil {
			return msg, err
		}
		if i == len(txpoolContent.QueuedContent)-1 {
			msg = msg + string(bytes)
		} else {
			msg = msg + string(bytes) + ",\n"
		}

	}
	msg = msg + "]}"
	return msg, nil
}

func TxPoolstatus() (model.TxPoolStatusResponse, error) {
	var response model.TxPoolStatusResponse
	data, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "txpool_status",
		"id":      config.ChainID,
		"params":  []interface{}{},
	})
	if err != nil {
		return response, err
	}

	resp, err := http.Post(config.BlockchainURL, constants.Response_format, strings.NewReader(string(data)))
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func Ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}
