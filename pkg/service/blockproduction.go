package service

import (
	"context"
	"math/big"
	"strconv"
	"time"

	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/config"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/email"
)

// check if blocks is getting mined
func (s *Service) checkBlockProduction() error {
	s.log.Info("checkBlockProduction start")

	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return err
	}

	currentBlock, err := s.ethClient.BlockNumber(context.Background())
	if err != nil {
		return err
	}
	s.log.Info("current block number: ", currentBlock)

	block, err := s.ethClient.BlockByNumber(context.Background(), new(big.Int).SetInt64(int64(currentBlock)))
	if err != nil {
		return err
	}

	if s.lastBlock == 0 {
		s.lastBlock = int(currentBlock)
		s.lastBlockMinedAt = int(block.Header().Time)
	} else if int(currentBlock) > s.lastBlock {
		s.lastBlock = int(currentBlock)
		s.lastBlockMinedAt = int(block.Header().Time)
	} else {
		if time.Now().Unix()-int64(s.lastBlockMinedAt) > int64(config.BlockProductionTime) {
			emailMessage := "Alert ! Block time exceeded " + strconv.Itoa(config.BlockProductionTime) + " seconds !\n\n" +
				"Last block was " + strconv.Itoa(s.lastBlock) + "\n" +
				"Last block was created at : " + time.Unix(int64(s.lastBlockMinedAt), 0).In(loc).String() + "\n" +
				"\n\nImportant : Number of block time exceeding emails skipped because of frequent emails is " + strconv.Itoa(s.blockProductionEmails.countOfEmailsSkipped)
			s.log.Infof(emailMessage)
			if time.Now().Unix()-s.blockProductionEmails.lastEmailsentAt > int64(config.EmailFrequency) {
				err := email.SendEmail(emailMessage)
				if err != nil {
					return err
				}
				s.blockProductionEmails.lastEmailsentAt = time.Now().Unix()
				s.blockProductionEmails.countOfEmailsSkipped = 0
			} else {
				s.log.Infof("Got frequent alerts of block time,%v email skipped", s.blockProductionEmails.countOfEmailsSkipped)
				s.blockProductionEmails.countOfEmailsSkipped = s.blockProductionEmails.countOfEmailsSkipped + 1
			}
		} else {
			s.log.Infof("Waiting for %v seconds for next block", config.BlockProductionTime)
		}
	}
	s.log.Info("checkBlockProduction end")
	return nil
}
