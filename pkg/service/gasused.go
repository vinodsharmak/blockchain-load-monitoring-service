package service

import (
	"context"
	"math/big"
	"strconv"

	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/config"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/email"
)

// Check if we have reached the max gas used and trigger email alerts
func (s *Service) checkForMaxGasUsed(startBlock int, endBlock int) error {
	maxGasUsedPerBlock, err := strconv.Atoi(config.MaxGasUsedPerBlock)
	if err != nil {
		return err
	}

	s.log.Infof("Calculating gas used by blocks between %v to %v ...", startBlock, endBlock)
	gasUsedLimitReached := false
	for i := startBlock; i <= endBlock; i++ {
		block, err := s.ethClient.BlockByNumber(context.Background(), new(big.Int).SetInt64(int64(i)))
		if err != nil {
			return err
		}
		if block.GasUsed() >= uint64(maxGasUsedPerBlock) {
			gasUsedLimitReached = true
		} else {
			gasUsedLimitReached = false
			break
		}
	}
	if gasUsedLimitReached {
		s.log.Infof("Gas used per block from %v to %v is higher than the set threshold of %v, please check the blockchain", startBlock, endBlock, maxGasUsedPerBlock)

		emaiMessage := "Alert ! \nBlockchain has reached its threshold for gas limit for range of blocks ! \n\n" +
			"Gas used from " + strconv.Itoa(startBlock) + " to " + strconv.Itoa(endBlock) +
			" has reached the Maximum gas used per block threshold of " + config.MaxGasUsedPerBlock +
			". Please check the blockchain. "

		err := email.SendEmail(emaiMessage)
		if err != nil {
			return err
		}
	}
	s.lastCheckedBlockForGasUSed = endBlock
	return nil
}

// check the current block and set the block range for gas used
func (s *Service) checkGasUsed() error {
	s.log.Info("checkGasUsed start")
	blockDifferenceForMaxGasUsed, err := strconv.Atoi(config.BlockDifferenceForMaxGasUsed)
	if err != nil {
		return err
	}

	currentBlock, err := s.ethClient.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	expectedBlock := 0
	if s.lastCheckedBlockForGasUSed == 0 {
		expectedBlock = int(currentBlock)
		s.lastCheckedBlockForGasUSed = expectedBlock - blockDifferenceForMaxGasUsed
	} else {
		expectedBlock = s.lastCheckedBlockForGasUSed + blockDifferenceForMaxGasUsed
	}

	if int(currentBlock) >= expectedBlock {
		err := s.checkForMaxGasUsed(s.lastCheckedBlockForGasUSed+1, expectedBlock)
		if err != nil {
			return err
		}
	}
	s.log.Info("ChecGasUsed end")
	return nil
}
