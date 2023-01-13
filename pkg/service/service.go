package service

import (
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

func (s *Service) StartTxCountMonitoring() error {
	err := s.checkTxLoad()
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
