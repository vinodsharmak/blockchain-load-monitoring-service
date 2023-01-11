package service

import (
	"encoding/json"
	"strconv"

	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/config"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/helpers"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/model"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func CheckPendingAndQueuedTxCount() error {
	log := config.Logger
	log.Info("CheckPendingAndQueuedTxCount start")

	maxTxPending, err := strconv.Atoi(config.MaxTxPending)
	if err != nil {
		return err
	}

	txpool_status, err := helpers.TxPoolstatus()
	if err != nil {
		return err
	}

	pendingHex := txpool_status.Result.Pending
	queuedHex := txpool_status.Result.Queued
	pending, err := hexutil.DecodeUint64(pendingHex)
	if err != nil {
		return err
	}
	queued, err := hexutil.DecodeUint64(queuedHex)
	if err != nil {
		return err
	}

	log.Infof("Total number of pending transactions is %v, Total number of queued transaction is %v", pending, queued)
	if pending >= uint64(maxTxPending) {
		log.Infof("Total number of pending transactions is %v, which is higher than the set threshold of %v, please check the blockchain.", pending, maxTxPending)
		txpool_content, err := helpers.TxPoolContent()
		if err != nil {
			return err
		}
		pending := txpool_content.Result.Pending
		var pendingTransactions []model.TxBody
		for key, value := range pending {
			var transaction model.TxBody
			transaction.From = key.String()
			for key2, value2 := range value {
				transaction.To = value2.To
				transaction.Nonce = strconv.Itoa(key2)
				transaction.Gas = value2.Gas
				transaction.Hash = value2.Hash
				pendingTransactions = append(pendingTransactions, transaction)
			}
		}
		pendingTransactionsJson, err := json.MarshalIndent(pendingTransactions, " ", " ")
		if err != nil {
			return err
		}
		pendingTransactionstring := string(pendingTransactionsJson)
		log.Infof(pendingTransactionstring)
		// TODO: send email
	}
	if queued > 0 {
		log.Infof("Total number of queued trnsactions is %v", queued)
		txpool_content, err := helpers.TxPoolContent()
		if err != nil {
			return err
		}
		queued := txpool_content.Result.Queued
		var queuedTransactions []model.TxBody
		for key, value := range queued {
			var transaction model.TxBody
			transaction.From = key.String()
			for key2, value2 := range value {
				transaction.To = value2.To
				transaction.Nonce = strconv.Itoa(key2)
				transaction.Gas = value2.Gas
				transaction.Hash = value2.Hash
				queuedTransactions = append(queuedTransactions, transaction)
			}
		}
		queuedTransactionsJson, err := json.MarshalIndent(queuedTransactions, " ", " ")
		if err != nil {
			return err
		}
		queuedTransactionstring := string(queuedTransactionsJson)
		log.Infof(queuedTransactionstring)
		// TODO: send email
	}

	log.Info("CheckPendingAndQueuedTxCount end")
	return nil
}
