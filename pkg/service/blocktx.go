package service

import (
	"context"
	"math/big"
	"strconv"

	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/config"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/model"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func (s *Service) Configure() error {
	var err error
	s.ethClient, err = ethclient.Dial(config.BlockchainURL)
	if err != nil {
		return err
	}
	s.log = config.Logger
	s.PendingTx = make(map[common.Address]map[uint64]model.TxBody)
	s.QueuedTx = make(map[common.Address]map[uint64]model.TxBody)
	return nil
}

// Check if we have reached the max tx load and trigger email alerts
func (s *Service) checkForMaxTxLoad(startBlock int, endBlock int) error {
	s.log.Info("CheckForMaxTxLoad start")

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
			// TODO: send email
		}
		totalTransactions = totalTransactions + int(transactionCount)
	}
	s.log.Infof("Total number of transaction between %v and %v is %v.", startBlock, endBlock, totalTransactions)

	if totalTransactions >= maxTxLoad {
		s.log.Infof("Transaction load is higher than the %v for %v blocks, Please check the blockchain.", maxTxLoad, config.BlockDifferenceForMaxTxLoad)
		// TODO: send email
	}

	s.lastCheckedBlock = endBlock
	s.log.Info("CheckForMaxTxLoad end")
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
		endBlock := s.lastCheckedBlock + blockDifferenceForMaxTxLoad - 1
		err := s.checkForMaxTxLoad(s.lastCheckedBlock, endBlock)
		if err != nil {
			return err
		}
	}

	return nil
}
