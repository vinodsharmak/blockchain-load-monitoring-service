package service

import (
	"context"
	"math/big"
	"strconv"

	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/config"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Check if we have reached the max tx load and trigger email alerts
func checkForMaxTxLoad(currentBlock int, oldBlock int) error {
	log := config.Logger
	log.Info("CheckForMaxTxLoad start")

	maxTxLoad, err := strconv.Atoi(config.MaxTxLoad)
	if err != nil {
		return err
	}
	maxTxPerBlock, err := strconv.Atoi(config.MaxTxPerBlock)
	if err != nil {
		return err
	}

	client, err := ethclient.Dial(config.BlockchainURL)
	if err != nil {
		return err
	}

	log.Infof("Calculating number of transactions between %v to %v ...", oldBlock, currentBlock)
	totalTransactions := 0
	for i := oldBlock; i <= currentBlock; i++ {
		block, err := client.BlockByNumber(context.Background(), new(big.Int).SetInt64(int64(i)))
		if err != nil {
			return err
		}
		transactionCount, err := client.TransactionCount(context.Background(), block.Hash())
		if err != nil {
			return err
		}
		if int(transactionCount) >= maxTxPerBlock {
			log.Infof("Transaction in %v was %v,which is higher than the tx/block threshold of %v Please check the blockchain.", i, int(transactionCount), maxTxPerBlock)
			// TODO: send email
		}
		totalTransactions = totalTransactions + int(transactionCount)
	}
	log.Infof("Total number of transaction between %v and %v is %v.", oldBlock, currentBlock, totalTransactions)

	if totalTransactions >= maxTxLoad {
		log.Infof("Transaction load is higher than the %v for %v blocks, Please check the blockchain.", maxTxLoad, config.BlockDifferenceForMaxTxLoad)
		// TODO: send email
	}

	log.Info("CheckForMaxTxLoad end")
	config.LastCheckedBlock = currentBlock
	return nil
}

func checkTxLoad() error {
	blockDifferenceForMaxTxLoad, err := strconv.Atoi(config.BlockDifferenceForMaxTxLoad)
	if err != nil {
		return err
	}

	client, err := ethclient.Dial(config.BlockchainURL)
	if err != nil {
		return err
	}
	currentBlock, err := client.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	expectedBlock := 0
	if config.LastCheckedBlock == 0 {
		expectedBlock = int(currentBlock) - blockDifferenceForMaxTxLoad
		config.LastCheckedBlock = expectedBlock
	} else {
		expectedBlock = config.LastCheckedBlock + blockDifferenceForMaxTxLoad
	}

	if int(currentBlock) >= expectedBlock {
		oldBlock := config.LastCheckedBlock
		newBlock := oldBlock + blockDifferenceForMaxTxLoad - 1
		err := checkForMaxTxLoad(newBlock, oldBlock)
		if err != nil {
			return err
		}
	}

	return nil
}
