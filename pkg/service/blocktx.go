package service

import (
	"context"
	"math/big"
	"strconv"

	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/config"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Service struct {
	ethClient                   *ethclient.Client
	maxTxLoad                   int
	maxTxPerBlock               int
	blockDifferenceForMaxTxLoad int
	lastCheckedBlock            int
}

func (s *Service) Configure() error {
	var err error
	s.ethClient, err = ethclient.Dial(config.BlockchainURL)
	if err != nil {
		return err
	}
	s.maxTxLoad, err = strconv.Atoi(config.MaxTxLoad)
	if err != nil {
		return err
	}
	s.maxTxPerBlock, err = strconv.Atoi(config.MaxTxPerBlock)
	if err != nil {
		return err
	}
	s.blockDifferenceForMaxTxLoad, err = strconv.Atoi(config.BlockDifferenceForMaxTxLoad)
	if err != nil {
		return err
	}

	return nil
}

// Check if we have reached the max tx load and trigger email alerts
func (s *Service) checkForMaxTxLoad(currentBlock int, oldBlock int) error {
	log := config.Logger
	log.Info("CheckForMaxTxLoad start")
	log.Infof("Calculating number of transactions between %v to %v ...", oldBlock, currentBlock)
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
		if int(transactionCount) >= s.maxTxPerBlock {
			log.Infof("Transaction in %v was %v,which is higher than the tx/block threshold of %v Please check the blockchain.", i, int(transactionCount), s.maxTxPerBlock)
			// TODO: send email
		}
		totalTransactions = totalTransactions + int(transactionCount)
	}
	log.Infof("Total number of transaction between %v and %v is %v.", oldBlock, currentBlock, totalTransactions)

	if totalTransactions >= s.maxTxLoad {
		log.Infof("Transaction load is higher than the %v for %v blocks, Please check the blockchain.", s.maxTxLoad, config.BlockDifferenceForMaxTxLoad)
		// TODO: send email
	}

	log.Info("CheckForMaxTxLoad end")
	s.lastCheckedBlock = currentBlock
	return nil
}

func (s *Service) checkTxLoad() error {
	currentBlock, err := s.ethClient.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	expectedBlock := 0
	if s.lastCheckedBlock == 0 {
		expectedBlock = int(currentBlock) - s.blockDifferenceForMaxTxLoad
		s.lastCheckedBlock = expectedBlock
	} else {
		expectedBlock = s.lastCheckedBlock + s.blockDifferenceForMaxTxLoad
	}

	if int(currentBlock) >= expectedBlock {
		oldBlock := s.lastCheckedBlock
		newBlock := oldBlock + s.blockDifferenceForMaxTxLoad - 1
		err := s.checkForMaxTxLoad(newBlock, oldBlock)
		if err != nil {
			return err
		}
	}

	return nil
}
