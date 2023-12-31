package service

import (
	"context"
	"math/big"
	"strconv"
	"time"

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
		emaiMessage := "Alert ! \nThreshold reached for block gas limit utilization! \n\n" +
			"Gas used from block - " + strconv.Itoa(startBlock) + " to block - " + strconv.Itoa(endBlock) +
			" has crossed the gas utilization per block threshold of " + config.MaxGasUsedPerBlock +
			"\n\nImportant : Number of gas utilization  alert emails skipped beacuse of frequency of emails is " + strconv.Itoa(s.GasUsedEmails.countOfEmailsSkipped)

		s.log.Infof(emaiMessage)

		if time.Now().Unix()-s.GasUsedEmails.lastEmailsentAt > int64(config.EmailFrequency) {
			err := email.SendEmail(emaiMessage)
			if err != nil {
				return err
			}
			s.GasUsedEmails.lastEmailsentAt = time.Now().Unix()
			s.GasUsedEmails.countOfEmailsSkipped = 0
		} else {
			s.log.Infof("Got frequent alerts of gas used,%v email skipped", s.GasUsedEmails.countOfEmailsSkipped)
			s.GasUsedEmails.countOfEmailsSkipped = s.GasUsedEmails.countOfEmailsSkipped + 1
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
	s.log.Info("CheckGasUsed end")
	return nil
}
