package service

import (
	"context"
	"math/big"
	"strconv"

	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/config"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/email"
)

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
			s.log.Infof("Transaction in %v was %v,which is higher than the tx/block threshold of %v Please check the blockchain.", i, int(transactionCount), maxTxPerBlock)

			emaiMessage := "Alert, Blockchain has reached its threshold for tx/block! \n\n" +
				"Maximum transaction threshold per block is " + config.MaxTxPerBlock + " and " + strconv.Itoa(i) + " has " +
				strconv.Itoa(int(transactionCount)) + "ttrnsactions. Please check the block."

			err := email.SendEmail(emaiMessage)
			if err != nil {
				s.log.Info("Error while sending Email")
				return err
			}
		}
		totalTransactions = totalTransactions + int(transactionCount)
	}
	s.log.Infof("Total number of transaction between %v and %v is %v.", startBlock, endBlock, totalTransactions)

	if totalTransactions >= maxTxLoad {
		s.log.Infof("Transaction load is higher than the %v for %v blocks, Please check the blockchain.", maxTxLoad, config.BlockDifferenceForMaxTxLoad)

		emaiMessage := "Alert, Blockchain has reached its threshold for tx/block for range of blocks ! \n\n" +
			"Maximum threshold per " + config.BlockDifferenceForMaxTxLoad + " blocks is " +
			config.MaxTxLoad + "\n" + "Number of transactions between " + strconv.Itoa(startBlock) +
			" and " + strconv.Itoa(endBlock) + " was " + strconv.Itoa(totalTransactions) + ". Please check the blocks."

		err := email.SendEmail(emaiMessage)
		if err != nil {
			s.log.Info("Error while sending Email")
			return err
		}
	}

	s.lastCheckedBlockForTxLoad = endBlock
	return nil
}

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
		expectedBlock = int(currentBlock) - blockDifferenceForMaxTxLoad
		s.lastCheckedBlockForTxLoad = expectedBlock
	} else {
		expectedBlock = s.lastCheckedBlockForTxLoad + blockDifferenceForMaxTxLoad
	}

	if int(currentBlock) >= expectedBlock {
		endBlock := s.lastCheckedBlockForTxLoad + blockDifferenceForMaxTxLoad - 1
		err := s.checkForMaxTxLoad(s.lastCheckedBlockForTxLoad, endBlock)
		if err != nil {
			return err
		}
	}
	s.log.Info("CheckTxLoad end")
	return nil
}
