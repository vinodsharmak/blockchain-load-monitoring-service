package service

import (
	"fmt"
	"time"

	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/config"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/email"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/model"
	"github.com/antigloss/go/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EmailDetails struct {
	lastEmailsentAt      int64
	countOfEmailsSkipped int
}

type Service struct {
	ethClient                  *ethclient.Client
	lastCheckedBlockForTxLoad  int
	lastCheckedBlockForGasUSed int
	log                        *logger.Logger
	PendingTx                  map[common.Address]map[uint64]model.TxBody
	QueuedTx                   map[common.Address]map[uint64]model.TxBody
	TxLoadEmails               EmailDetails
	GasUsedEmails              EmailDetails
	PendingAndQueuedTxEmails   EmailDetails
	TxPoolStuckEmails          EmailDetails
	blockProductionEmails      EmailDetails
	lastBlock                  int
	lastBlockMinedAt           int
	failures                   int
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
		//check block poduction
		err := s.checkBlockProduction()
		if err != nil {
			s.failures++
			s.log.Errorf("checkBlockProduction: %v", err)
			if s.failures > 3 {
				err = email.SendEmail("SERVICE DOWN!!\nError encountered in service: " + err.Error())
				if err != nil {
					s.log.Errorf("error in sendEmail: %v", err)
					return err
				}
				return fmt.Errorf("checkTxLoad: %s", err)
			}
		} else {
			s.failures = 0
		}

		time.Sleep(time.Second * time.Duration(config.TimeIntervalForSubService))

		//check transaction load on blockchain
		err = s.checkTxLoad()
		if err != nil {
			s.failures++
			s.log.Errorf("checkTxLoad: %v", err)
			if s.failures > 3 {
				err = email.SendEmail("SERVICE DOWN!!\nError encountered in service: " + err.Error())
				if err != nil {
					s.log.Errorf("error in sendEmail: %v", err)
					return err
				}
				return fmt.Errorf("checkTxLoad: %s", err)
			}
		} else {
			s.failures = 0
		}

		time.Sleep(time.Second * time.Duration(config.TimeIntervalForSubService))

		//check block gaslimit usage
		err = s.checkGasUsed()
		if err != nil {
			s.failures++
			s.log.Errorf("checkGasUsed: %v", err)
			if s.failures > 3 {
				err = email.SendEmail("SERVICE DOWN!!\nError encountered in service: " + err.Error())
				if err != nil {
					s.log.Errorf("error in sendEmail: %v", err)
					return err
				}
				return fmt.Errorf("checkGasUsed: %s", err)
			}
		} else {
			s.failures = 0
		}

		time.Sleep(time.Second * time.Duration(config.TimeIntervalForSubService))

		//check pending and queued txpool count
		err = s.checkPendingAndQueuedTxCount()
		if err != nil {
			s.failures++
			s.log.Errorf("checkPendingAndQueuedTxCount: %v", err)
			if s.failures > 3 {
				err = email.SendEmail("SERVICE DOWN!!\nError encountered in service: " + err.Error())
				if err != nil {
					s.log.Errorf("error in sendEmail: %v", err)
					return err
				}
				return fmt.Errorf("checkPendingAndQueuedTxCount: %s", err)
			}
		} else {
			s.failures = 0
		}

		time.Sleep(time.Second * time.Duration(config.TimeIntervalForSubService))

		//check if tx stuck in txpool in pending and queued
		err = s.txPoolStuck()
		if err != nil {
			s.failures++
			s.log.Errorf("txPoolStuck: %v", err)
			if s.failures > 3 {
				err = email.SendEmail("SERVICE DOWN!!\nError encountered in service: " + err.Error())
				if err != nil {
					s.log.Errorf("error in sendEmail: %v", err)
					return err
				}
				return fmt.Errorf("txPoolStuck: %s", err)
			}
		} else {
			s.failures = 0
		}

		time.Sleep(time.Second * time.Duration(config.TimeIntervalForSubService))

		s.log.Info("-----------------------------------------------CYCLE-END-----------------------------------------------")
	}
}
