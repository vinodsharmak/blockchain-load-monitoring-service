package service

import (
	"fmt"
	"time"

	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/config"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/model"
	"github.com/antigloss/go/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Service struct {
	ethClient                  *ethclient.Client
	lastCheckedBlockForTxLoad  int
	lastCheckedBlockForGasUSed int
	log                        *logger.Logger
	PendingTx                  map[common.Address]map[uint64]model.TxBody
	QueuedTx                   map[common.Address]map[uint64]model.TxBody
}

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

func (s *Service) BlockchainMonitoringService() error {
	s.log.Info("starting blockchain monitoring service")

	for {
		s.log.Info("-----------------------------------------------CYCLE-START-----------------------------------------------")
		//check transaction load on blockchain
		err := s.checkTxLoad()
		if err != nil {
			return fmt.Errorf("checkTxLoad: %s", err)
		}

		time.Sleep(time.Second * time.Duration(config.TimeIntervalForSubService))

		//check block gaslimit usage
		err = s.checkGasUsed()
		if err != nil {
			return fmt.Errorf("checkGasUsed: %s", err)

		}

		time.Sleep(time.Second * time.Duration(config.TimeIntervalForSubService))

		//check pending and queued txpool count
		err = s.checkPendingAndQueuedTxCount()
		if err != nil {
			return fmt.Errorf("checkPendingAndQueuedTxCount: %s", err)
		}

		time.Sleep(time.Second * time.Duration(config.TimeIntervalForSubService))

		//check if tx stuck in txpool in pending and queued
		err = s.txPoolStuck()
		if err != nil {
			return fmt.Errorf("txPoolStuck: %s", err)
		}

		time.Sleep(time.Second * time.Duration(config.TimeIntervalForSubService))

		s.log.Info("-----------------------------------------------CYCLE-END-----------------------------------------------")
	}
}
