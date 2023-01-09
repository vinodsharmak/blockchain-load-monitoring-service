package service

import (
	"context"
	"math/big"
	"strconv"

	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/config"
	"github.com/antigloss/go/logger"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Service struct {
	ethClient        *ethclient.Client
	lastCheckedBlock int
	log              *logger.Logger
}

func (s *Service) Configure() error {
	var err error
	s.ethClient, err = ethclient.Dial(config.BlockchainURL)
	if err != nil {
		return err
	}
	s.log = config.Logger
	return nil
}

// Check if we have reached the max tx load and trigger email alerts
func (s *Service) checkForMaxTxLoad(currentBlock int, oldBlock int) error {
	s.log.Info("CheckForMaxTxLoad start")

	maxTxLoad, err := strconv.Atoi(config.MaxTxLoad)
	if err != nil {
		return err
	}
	maxTxPerBlock, err := strconv.Atoi(config.MaxTxPerBlock)
	if err != nil {
		return err
	}

	s.log.Infof("Calculating number of transactions between %v to %v ...", oldBlock, currentBlock)
	totalTransactions := 0
	for i := oldBlock; i <= currentBlock; i++ {
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
			// TODO: send email
		}
		totalTransactions = totalTransactions + int(transactionCount)
	}
	s.log.Infof("Total number of transaction between %v and %v is %v.", oldBlock, currentBlock, totalTransactions)

	if totalTransactions >= maxTxLoad {
		s.log.Infof("Transaction load is higher than the %v for %v blocks, Please check the blockchain.", maxTxLoad, config.BlockDifferenceForMaxTxLoad)
		// TODO: send email
	}

	s.log.Info("CheckForMaxTxLoad end")
	s.lastCheckedBlock = currentBlock
	return nil
}

// check the current block and set the block range for tx load test
func (s *Service) checkTxLoad() error {
	blockDifferenceForMaxTxLoad, err := strconv.Atoi(config.BlockDifferenceForMaxTxLoad)
	if err != nil {
		return err
	}

	currentBlock, err := s.ethClient.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	expectedBlock := 0
	if s.lastCheckedBlock == 0 {
		expectedBlock = int(currentBlock) - blockDifferenceForMaxTxLoad
		s.lastCheckedBlock = expectedBlock
	} else {
		expectedBlock = s.lastCheckedBlock + blockDifferenceForMaxTxLoad
	}

	if int(currentBlock) >= expectedBlock {
		oldBlock := s.lastCheckedBlock
		newBlock := oldBlock + blockDifferenceForMaxTxLoad - 1
		err := s.checkForMaxTxLoad(newBlock, oldBlock)
		if err != nil {
			return err
		}
	}

	return nil
}
