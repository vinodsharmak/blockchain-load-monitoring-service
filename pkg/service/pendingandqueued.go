package service

import (
	"strconv"

	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/config"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/email"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/helpers"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/model"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var txPoolCountMail model.TxPoolContentMail

func (s *Service) checkPendingAndQueuedTxCount() error {
	s.log.Info("CheckPendingAndQueuedTxCount start")
	txPoolCountMail = model.TxPoolContentMail{}

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

	txPoolCountMail.PendingCount = int(pendingTxCount)
	txPoolCountMail.QueuedCount = int(queuedTxCount)

	s.log.Infof("Total number of pending transactions is %v, Total number of queued transaction is %v", pendingTxCount, queuedTxCount)
	if pendingTxCount >= uint64(maxTxPending) || queuedTxCount > 0 {
		err := getPendingTransactionDetails()
		if err != nil {
			return err
		}
		err = getQueuedTransactionDetails()
		if err != nil {
			return err
		}

		emailMessage := "Alert ! \nThreshold reached for transaction pool! \n" +
			"Please find the transactions detail below : \n"
		if txPoolCountMail.PendingCount >= maxTxPending {
			emailMessage = emailMessage + "Number of pending transactions is higher than the threshold of  " +
				config.MaxTxPending + "!\n"
		}
		if txPoolCountMail.QueuedCount > 0 {
			emailMessage = emailMessage + "Blockchain have queued transactions in the pool!\n"

		}
		txpoolContentString, err := helpers.PrepareEmailBodyForTxPoolContent(txPoolCountMail)
		if err != nil {
			return err
		}
		emailMessage = emailMessage + txpoolContentString
		s.log.Infof(emailMessage)
		err = email.SendEmail(emailMessage)
		if err != nil {
			return err
		}
	}
	s.log.Info("checkPendingAndQueuedTxCount end")

	return nil
}

func getPendingTransactionDetails() error {
	txpool_content, err := helpers.TxPoolContent()
	if err != nil {
		return err
	}
	pending := txpool_content.Result.Pending
	for _, transactions := range pending {
		for nonce, txBody := range transactions {
			txBody.Nonce = strconv.Itoa(int(nonce))
			txPoolCountMail.PendingContent = append(txPoolCountMail.PendingContent, txBody)
		}
	}
	return err
}

func getQueuedTransactionDetails() error {
	txpool_content, err := helpers.TxPoolContent()
	if err != nil {
		return err
	}
	queued := txpool_content.Result.Queued
	for _, transactions := range queued {
		for nonce, txBody := range transactions {
			txBody.Nonce = strconv.Itoa(int(nonce))
			txPoolCountMail.QueuedContent = append(txPoolCountMail.QueuedContent, txBody)
		}
	}
	return err
}
