package service

import (
	"context"
	"math/big"
	"strconv"

	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/config"
	"github.com/antigloss/go/logger"
	"github.com/ethereum/go-ethereum/ethclient"
)

type GasUsedService struct {
	ethClient        *ethclient.Client
	lastCheckedBlock int
	log              *logger.Logger
}

func (s *GasUsedService) Configure() error {
	var err error
	s.ethClient, err = ethclient.Dial(config.BlockchainURL)
	if err != nil {
		return err
	}
	s.log = config.Logger
	return nil
}

// Check if we have reached the max gas used and trigger email alerts
func (s *GasUsedService) checkForMaxGasUsed(startBlock int, endBlock int) error {
	s.log.Info("CheckForMaxTxLoad start")

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
		// TODO: Send email
	}
	s.lastCheckedBlock = endBlock
	s.log.Info("CheckForMaxTxLoad end")
	return nil
}

// check the current block and set the block range for gas used
func (s *GasUsedService) checkGasUsed() error {
	blockDifferenceForMaxGasUsed, err := strconv.Atoi(config.BlockDifferenceForMaxGasUsed)
	if err != nil {
		return err
	}

	currentBlock, err := s.ethClient.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	expectedBlock := 0
	if s.lastCheckedBlock == 0 {
		expectedBlock = int(currentBlock) - blockDifferenceForMaxGasUsed
		s.lastCheckedBlock = expectedBlock
	} else {
		expectedBlock = s.lastCheckedBlock + blockDifferenceForMaxGasUsed
	}

	if int(currentBlock) >= expectedBlock {
		endBlock := s.lastCheckedBlock + blockDifferenceForMaxGasUsed - 1
		err := s.checkForMaxGasUsed(s.lastCheckedBlock, endBlock)
		if err != nil {
			return err
		}
	}

	return nil
}
