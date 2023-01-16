package service

import (
	"encoding/json"
	"fmt"
	"strconv"

	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/config"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/email"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/helpers"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/model"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (s *Service) checkPendingAndQueuedTxCount() error {
	s.log.Info("CheckPendingAndQueuedTxCount start")

	maxTxPending, err := strconv.Atoi(config.MaxTxPending)
	if err != nil {
		return err
	}

	txpoolStatus, err := helpers.TxPoolstatus()
	if err != nil {
		return err
	}

	pendingTxCount, err := hexutil.DecodeUint64(txpoolStatus.Result.Pending)
	if err != nil {
		return err
	}
	queuedTxCount, err := hexutil.DecodeUint64(txpoolStatus.Result.Queued)
	if err != nil {
		return err
	}
	s.log.Infof("Total number of pending transactions is %v, Total number of queued transaction is %v", pendingTxCount, queuedTxCount)
	if pendingTxCount >= uint64(maxTxPending) || queuedTxCount > 0 {
		pendingTransactionDetails, err := getPendingTransactionDetails()
		if err != nil {
			return err
		}
		queuedTransactionDetails, err := getQueuedTransactionDetails()
		if err != nil {
			return err
		}
		pendingTransactionsBytes, err := json.MarshalIndent(pendingTransactionDetails, " ", "")
		if err != nil {
			return err
		}
		pendingTransactionString := string(pendingTransactionsBytes)
		queuedTransactionsBytes, err := json.MarshalIndent(queuedTransactionDetails, " ", " ")
		if err != nil {
			return err
		}
		queuedTransactionsString := string(queuedTransactionsBytes)
		emailMessage := "Alert from pending and queued transaction count check ! \n" +
			"Please find the transaction pool details below : \n "
		if len(pendingTransactionDetails) >= maxTxPending {
			emailMessage = emailMessage + "Pending transaction count is higher than the threshold of " +
				config.MaxTxPending + "! \n "
			emailMessage = emailMessage + "Pending Transaction details : \n " + pendingTransactionString + "\n"
		}
		if len(queuedTransactionDetails) > 0 {
			emailMessage = "Blockchain have queued transactions in the pool ! "
			emailMessage = emailMessage + "Queued Transaction details : \n " + queuedTransactionsString + "\n"
		}
		s.log.Infof(emailMessage)
		err = email.SendEmail(emailMessage)
		if err != nil {
			s.log.Info("Error while sending Email")
			return err
		}
	}
	s.log.Info("CheckPendingAndQueuedTxCount end")

	return nil
}

func getPendingTransactionDetails() ([]model.TxBody, error) {
	var pendingTransactions []model.TxBody
	txpool_content, err := helpers.TxPoolContent()
	if err != nil {
		return pendingTransactions, err
	}
	pending := txpool_content.Result.Pending
	for fromAddress, transactions := range pending {
		var transaction model.TxBody
		transaction.From = fromAddress
		for nonce, txBody := range transactions {
			transaction.To = txBody.To
			transaction.Nonce = fmt.Sprintf("%v", nonce)
			transaction.Gas = txBody.Gas
			transaction.Hash = txBody.Hash
			pendingTransactions = append(pendingTransactions, transaction)
		}
	}
	return pendingTransactions, err
}

func getQueuedTransactionDetails() ([]model.TxBody, error) {
	var queuedTransactions []model.TxBody
	txpool_content, err := helpers.TxPoolContent()
	if err != nil {
		return queuedTransactions, err
	}
	queued := txpool_content.Result.Queued
	for fromAddress, transactions := range queued {
		var transaction model.TxBody
		transaction.From = fromAddress
		for nonce, txBody := range transactions {
			transaction.To = txBody.To
			transaction.Nonce = fmt.Sprintf("%v", nonce)
			transaction.Gas = txBody.Gas
			transaction.Hash = txBody.Hash
			queuedTransactions = append(queuedTransactions, transaction)
		}
	}
	return queuedTransactions, err
}
