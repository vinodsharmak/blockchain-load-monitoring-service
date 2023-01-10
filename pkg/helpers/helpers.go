package helpers

import (
	"encoding/json"
	"net/http"
	"strings"

	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/config"
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

	resp, err := http.Post(config.BlockchainURL, "application/json", strings.NewReader(string(data)))
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
