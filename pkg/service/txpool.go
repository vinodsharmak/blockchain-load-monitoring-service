package service

import (
	"fmt"
	"time"

	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/config"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/helpers"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/model"
	"github.com/ethereum/go-ethereum/common"
)

var txPoolContentStuckMail model.TxPoolContentStuckMail

func (s *Service) txPoolStuck() error {
	txPoolContentStuckMail = model.TxPoolContentStuckMail{}
	res, err := helpers.TxPoolContent()
	if err != nil {
		return err
	}

	//pending
	if len(res.Result.Pending) != 0 {
		s.removeConfirmedPending(res.Result.Pending)
		s.updatePendingContent(res.Result.Pending)
	} else {
		s.PendingTx = make(map[common.Address]map[uint64]model.TxBody)
	}

	//queued
	if len(res.Result.Queued) != 0 {
		s.removeConfirmedQueued(res.Result.Queued)
		s.updateQueuedContent(res.Result.Queued)
	} else {
		s.QueuedTx = make(map[common.Address]map[uint64]model.TxBody)
	}

	// mailContentBytes, err := json.MarshalIndent(txPoolContentStuckMail, "\n", "\n")
	// if err != nil {
	// 	return err
	// }
	// json.Marshal(txPoolContentStuckMail)
	mailContentString, err := helpers.PrepareEmailBodyForTxPoolContent(txPoolContentStuckMail)
	if err != nil {
		return err
	}
	fmt.Println("TXPOOL STUCK MAIL BODY:\n", mailContentString)
	// fmt.Println("TXPOOL LOCAL:", s.PendingTx)

	return nil
}

func (s *Service) updatePendingContent(pending map[common.Address]map[uint64]model.TxBody) {
	for sender, transactions := range pending {
		if s.checkIfSenderAvailableInPending(sender) {
			for nonce, txBody := range transactions {
				if s.checkIfNonceAvailableInPending(sender, nonce) {
					timeElapsedInSeconds := time.Now().Unix() - s.PendingTx[sender][nonce].FoundAtEpochTime
					if timeElapsedInSeconds >= int64(config.TxpoolTimeLimit) {
						txPoolContentStuckMail.PendingCount++
						txPoolContentStuckMail.PendingContent = append(txPoolContentStuckMail.PendingContent, s.PendingTx[sender][nonce])
					}
				} else {
					txBody.FoundAtEpochTime = time.Now().Unix()
					s.PendingTx[sender][nonce] = txBody
				}
			}
		} else {
			for nonce, txBody := range transactions {
				txBody.FoundAtEpochTime = time.Now().Unix()
				s.PendingTx[sender] = make(map[uint64]model.TxBody)
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
						txPoolContentStuckMail.QueuedCount++
						txPoolContentStuckMail.QueuedContent = append(txPoolContentStuckMail.QueuedContent, s.QueuedTx[sender][nonce])
					}
				} else {
					txBody.FoundAtEpochTime = time.Now().Unix()
					s.QueuedTx[sender][nonce] = txBody
				}
			}
		} else {
			for nonce, txBody := range transactions {
				txBody.FoundAtEpochTime = time.Now().Unix()
				s.QueuedTx[sender] = make(map[uint64]model.TxBody)
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
			for nonce, _ := range transactions {
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
			for nonce, _ := range transactions {
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