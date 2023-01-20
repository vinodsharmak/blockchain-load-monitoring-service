package service

import (
	"context"
	"encoding/json"
	"math/big"
	"strconv"
	"time"

	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/config"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/email"
)

// check the current block and set the block range for tx load test
func (s *Service) checkTxLoad() error {
	s.log.Info("checkTxLoad start")
	blockDifferenceForMaxTxLoad, err := strconv.Atoi(config.BlockDifferenceForMaxTxLoad)
	if err != nil {
		return err
	}

	currentBlock, err := s.ethClient.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	expectedBlock := 0
	if s.lastCheckedBlockForTxLoad == 0 {
		expectedBlock = int(currentBlock)
		s.lastCheckedBlockForTxLoad = expectedBlock - blockDifferenceForMaxTxLoad
	} else {
		expectedBlock = s.lastCheckedBlockForTxLoad + blockDifferenceForMaxTxLoad
	}

	if int(currentBlock) >= expectedBlock {
		err := s.checkForMaxTxLoad(s.lastCheckedBlockForTxLoad+1, expectedBlock)
		if err != nil {
			return err
		}
	}
	s.log.Info("CheckTxLoad end")
	return nil
}

// Check if we have reached the max tx load and trigger email alerts
func (s *Service) checkForMaxTxLoad(startBlock int, endBlock int) error {
	maxTxLoad, err := strconv.Atoi(config.MaxTxLoad)
	if err != nil {
		return err
	}
	maxTxPerBlock, err := strconv.Atoi(config.MaxTxPerBlock)
	if err != nil {
		return err
	}

	higherTxLoadBlocks := make(map[int]uint)

	s.log.Infof("Calculating number of transactions between %v to %v ...", startBlock, endBlock)
	totalTransactions := 0
	for i := startBlock; i <= endBlock; i++ {
		block, err := s.ethClient.BlockByNumber(context.Background(), new(big.Int).SetInt64(int64(i)))
		if err != nil {
			return err
		}
		transactionCount, err := s.ethClient.TransactionCount(context.Background(), block.Hash())
		if err != nil {
			return err
		}
		if int(transactionCount) >= maxTxPerBlock {
			s.log.Infof("Transaction in %v was %v,which is higher than the tx/block threshold of %v", i, transactionCount, maxTxPerBlock)
			higherTxLoadBlocks[i] = transactionCount
		}
		totalTransactions = totalTransactions + int(transactionCount)
	}
	s.log.Infof("Total number of transaction between %v and %v is %v.", startBlock, endBlock, totalTransactions)

	if totalTransactions >= maxTxLoad {
		s.log.Infof("Transaction load is higher than the %v for %v blocks.", maxTxLoad, config.BlockDifferenceForMaxTxLoad)
		emaiMessage := "Alert !\nThreshold reached for total transactions within a block range! \n\n" +
			"Maximum threshold per " + config.BlockDifferenceForMaxTxLoad + " blocks is " +
			config.MaxTxLoad + "\n" + "Number of transactions between " + strconv.Itoa(startBlock) +
			" and " + strconv.Itoa(endBlock) + " was " + strconv.Itoa(totalTransactions) +
			"\n\nImportant : Number of  transaction load alert emails skipped beacuse of frequency of emails is " + strconv.Itoa(s.TxLoadEmails.countOfEmailsSkipped)

		s.log.Infof(emaiMessage)

		if time.Now().Unix()-s.TxLoadEmails.lastEmailsentAt > int64(config.EmailFrequency) {
			err := email.SendEmail(emaiMessage)
			if err != nil {
				return err
			}
			s.TxLoadEmails.lastEmailsentAt = time.Now().Unix()
			s.TxLoadEmails.countOfEmailsSkipped = 0
		} else {
			s.log.Infof("Got frequent alerts of tx load ,%v email skipped", s.TxLoadEmails.countOfEmailsSkipped)
			s.TxLoadEmails.countOfEmailsSkipped = s.TxLoadEmails.countOfEmailsSkipped + 1
		}
	} else {
		if len(higherTxLoadBlocks) > 0 {
			higherTxLoadBlocksBytes, err := json.Marshal(higherTxLoadBlocks)
			if err != nil {
				return err
			}
			emaiMessage := "Alert ! \n Threshold reached for transactions per block! \n\n" +
				"Threshold of transactions per block is " + config.MaxTxPerBlock + "\n" +
				"These blocks has passed the threshold of transactions count per block : \n" +
				"Given detail is in format : {Block Number: Transactions count} \n" + string(higherTxLoadBlocksBytes) +
				"\n\nImportant : Number of transaction load alert emails skipped beacuse of frequency of emails is " + strconv.Itoa(s.TxLoadEmails.countOfEmailsSkipped)

			s.log.Infof(emaiMessage)

			if time.Now().Unix()-s.TxLoadEmails.lastEmailsentAt > int64(config.EmailFrequency) {
				err := email.SendEmail(emaiMessage)
				if err != nil {
					return err
				}
				s.TxLoadEmails.lastEmailsentAt = time.Now().Unix()
				s.TxLoadEmails.countOfEmailsSkipped = 0
			} else {
				s.log.Infof("Got frequent alerts of tx load ,%v email skipped", s.TxLoadEmails.countOfEmailsSkipped)
				s.TxLoadEmails.countOfEmailsSkipped = s.TxLoadEmails.countOfEmailsSkipped + 1
			}
		}
	}

	s.lastCheckedBlockForTxLoad = endBlock
	return nil
}
