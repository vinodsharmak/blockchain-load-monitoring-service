package service

import (
	"time"

	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/config"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/model"
	"github.com/antigloss/go/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Service struct {
	ethClient        *ethclient.Client
	lastCheckedBlock int
	log              *logger.Logger
	PendingTx        map[common.Address]map[uint64]model.TxBody
	QueuedTx         map[common.Address]map[uint64]model.TxBody
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

func (s *Service) StartTxCountMonitoring() error {
	err := s.txPoolStuck()
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 10)

	err = s.checkTxLoad()
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) StartPendingAndQueuedTxMonitoring() error {
	err := s.CheckPendingAndQueuedTxCount()
	if err != nil {
		return err
	}
	return nil
}
