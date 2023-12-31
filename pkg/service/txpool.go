package service

import (
	"strconv"
	"time"

	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/config"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/email"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/helpers"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/model"
	"github.com/ethereum/go-ethereum/common"
)

var txPoolContentStuckMail model.TxPoolContentMail

func (s *Service) txPoolStuck() error {
	s.log.Infof("txPoolStuck started")
	txPoolContentStuckMail = model.TxPoolContentMail{}
	res, err := helpers.TxPoolContent()
	if err != nil {
		return err
	}

	//pending
	s.log.Infof("txPoolStuck checking pending")
	if len(res.Result.Pending) != 0 {
		s.removeConfirmedPending(res.Result.Pending)
		s.updatePendingContent(res.Result.Pending)
	} else {
		s.PendingTx = make(map[common.Address]map[uint64]model.TxBody)
	}

	s.log.Infof("txPoolStuck checking queued")
	//queued
	if len(res.Result.Queued) != 0 {
		s.removeConfirmedQueued(res.Result.Queued)
		s.updateQueuedContent(res.Result.Queued)
	} else {
		s.QueuedTx = make(map[common.Address]map[uint64]model.TxBody)
	}

	mailContentString, err := helpers.PrepareEmailBodyForTxPoolContent(txPoolContentStuckMail)
	if err != nil {
		return err
	}

	s.log.Info("txPoolStuck result:\n", mailContentString)
	if txPoolContentStuckMail.PendingCount > 0 || txPoolContentStuckMail.QueuedCount > 0 {
		emailMessage := "Alert !\nTransactions are stuck inside the transaction pool (pending/queued), please find the details below :\n" + mailContentString +
			"\n\nImportant : Number of transaction pool stuck alert emails skipped beacuse of frequent emails is " + strconv.Itoa(s.TxPoolStuckEmails.countOfEmailsSkipped)

		s.log.Infof(emailMessage)

		if time.Now().Unix()-s.TxPoolStuckEmails.lastEmailsentAt > int64(config.EmailFrequency) {
			err := email.SendEmail(emailMessage)
			if err != nil {
				return err
			}
			s.TxPoolStuckEmails.lastEmailsentAt = time.Now().Unix()
			s.TxPoolStuckEmails.countOfEmailsSkipped = 0
		} else {
			s.log.Infof("Got frequent alerts of tx pool stuck transaction,%v email skipped", s.TxPoolStuckEmails.countOfEmailsSkipped)
			s.TxPoolStuckEmails.countOfEmailsSkipped = s.TxPoolStuckEmails.countOfEmailsSkipped + 1
		}
	}

	s.log.Infof("txPoolStuck end")

	return nil
}

func (s *Service) updatePendingContent(pending map[common.Address]map[uint64]model.TxBody) {
	for sender, transactions := range pending {
		if s.checkIfSenderAvailableInPending(sender) {
			for nonce, txBody := range transactions {
				if s.checkIfNonceAvailableInPending(sender, nonce) {
					timeElapsedInSeconds := time.Now().Unix() - s.PendingTx[sender][nonce].FoundAtEpochTime
					if timeElapsedInSeconds >= int64(config.TxpoolTimeLimit) {
						s.log.Infof("found tx stuck in pending: %s", s.PendingTx[sender][nonce])
						txPoolContentStuckMail.PendingCount++
						txPoolContentStuckMail.PendingContent = append(txPoolContentStuckMail.PendingContent, s.PendingTx[sender][nonce])
					}
				} else {
					txBody.FoundAtEpochTime = time.Now().Unix()
					txBody.Nonce = strconv.Itoa(int(nonce))
					s.PendingTx[sender][nonce] = txBody
				}
			}
		} else {
			for nonce, txBody := range transactions {
				txBody.FoundAtEpochTime = time.Now().Unix()
				s.PendingTx[sender] = make(map[uint64]model.TxBody)
				txBody.Nonce = strconv.Itoa(int(nonce))
				s.PendingTx[sender][nonce] = txBody
			}
		}

	}
}

func (s *Service) updateQueuedContent(queued map[common.Address]map[uint64]model.TxBody) {
	for sender, transactions := range queued {
		if s.checkIfSenderAvailableInQueued(sender) {
			for nonce, txBody := range transactions {
				if s.checkIfNonceAvailableInQueued(sender, nonce) {
					timeElapsedInSeconds := time.Now().Unix() - s.PendingTx[sender][nonce].FoundAtEpochTime
					if timeElapsedInSeconds >= int64(config.TxpoolTimeLimit) {
						s.log.Infof("found tx stuck in queued: %s", s.QueuedTx[sender][nonce])
						txPoolContentStuckMail.QueuedCount++
						txBody.Nonce = strconv.Itoa(int(nonce))
						txPoolContentStuckMail.QueuedContent = append(txPoolContentStuckMail.QueuedContent, s.QueuedTx[sender][nonce])
					}
				} else {
					txBody.FoundAtEpochTime = time.Now().Unix()
					txBody.Nonce = strconv.Itoa(int(nonce))
					s.QueuedTx[sender][nonce] = txBody
				}
			}
		} else {
			for nonce, txBody := range transactions {
				txBody.FoundAtEpochTime = time.Now().Unix()
				s.QueuedTx[sender] = make(map[uint64]model.TxBody)
				txBody.Nonce = strconv.Itoa(int(nonce))
				s.QueuedTx[sender][nonce] = txBody
			}
		}

	}
}

func (s *Service) removeConfirmedPending(pending map[common.Address]map[uint64]model.TxBody) {
	for sender, transactions := range s.PendingTx {
		_, ok := pending[sender]
		if !ok {
			delete(s.PendingTx, sender)
		} else {
			for nonce := range transactions {
				_, ok := pending[sender][nonce]
				if !ok {
					delete(s.PendingTx[sender], nonce)
				}
			}
		}
	}
}

func (s *Service) removeConfirmedQueued(queued map[common.Address]map[uint64]model.TxBody) {
	for sender, transactions := range s.QueuedTx {
		_, ok := queued[sender]
		if !ok {
			delete(s.QueuedTx, sender)
		} else {
			for nonce := range transactions {
				_, ok := queued[sender][nonce]
				if !ok {
					delete(s.QueuedTx[sender], nonce)
				}
			}
		}
	}
}

func (s *Service) checkIfSenderAvailableInPending(address common.Address) bool {
	_, available := s.PendingTx[address]
	return available
}

func (s *Service) checkIfSenderAvailableInQueued(address common.Address) bool {
	_, available := s.QueuedTx[address]
	return available
}

func (s *Service) checkIfNonceAvailableInPending(address common.Address, nonce uint64) bool {
	_, ok := s.PendingTx[address][nonce]
	return ok
}

func (s *Service) checkIfNonceAvailableInQueued(address common.Address, nonce uint64) bool {
	_, ok := s.QueuedTx[address][nonce]
	return ok
}
