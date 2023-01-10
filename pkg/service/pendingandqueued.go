package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/config"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/constants"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func CheckPendingAndQueuedTxCount() error {
	log := config.Logger
	log.Info("CheckPendingAndQueuedTxCount start")

	maxTxPending, err := strconv.Atoi(config.MaxTxPending)
	if err != nil {
		return err
	}

	data, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "txpool_status",
		"id":      config.ChainID,
		"params":  []interface{}{},
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(config.BlockchainURL,
		constants.Response_format, strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	result := make(map[string]interface{})
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	pendingHex := result["result"].(map[string]interface{})["pending"]
	queuedHex := result["result"].(map[string]interface{})["queued"]
	pending, err := hexutil.DecodeUint64(fmt.Sprintf("%s", pendingHex))
	if err != nil {
		return err
	}
	queued, err := hexutil.DecodeUint64(fmt.Sprintf("%s", queuedHex))
	if err != nil {
		return err
	}

	log.Infof("Total number of pending transactions is %v, Total number of queued transaction is %v", pending, queued)

	if err != nil {
		log.Info("Error while gettting details", err)
	}

	if pending >= uint64(maxTxPending) || queued > 0 {
		log.Infof("Total number of pending transactions is %v, which is higher than the set threshold of %v, please check the blockchain.", pending, maxTxPending)
		// TODO: send email
	}

	log.Info("CheckPendingAndQueuedTxCount end")
	return nil
}
