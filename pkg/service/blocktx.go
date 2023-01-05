package service

import (
	"context"
	"math/big"
	"strconv"

	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/config"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/constants"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/email"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Check if we have reached the max tx load and trigger email alerts
func CheckForMaxTxLoad() error {
	logging := config.Logger
	logging.Info("CheckForMaxTxLoad start")

	maxTxLoad, _ := strconv.Atoi(config.MaxTxLoad)
	blockDifferenceForMaxTxLoad, _ := strconv.ParseInt(config.BlockDifferenceForMaxTxLoad, 10, 64)

	cl, err := ethclient.Dial(config.BlockchainURL)
	if err != nil {
		logging.Errorf(" Error while configuring eth client : %v", err.Error())
		return err
	}

	currentBlock, err := cl.HeaderByNumber(context.Background(), nil)
	if err != nil {
		logging.Errorf(" Error while getting latest block : %v", err.Error())
		return err
	}
	endBlock := currentBlock.Number.Int64()
	startBlock := endBlock - blockDifferenceForMaxTxLoad

	logging.Infof("Calculating number of transactions between %v to %v ...", startBlock, endBlock)
	totalTransactions := 0
	for i := startBlock; i <= endBlock; i++ {
		block, err := cl.BlockByNumber(context.Background(), new(big.Int).SetInt64(i))
		if err != nil {
			logging.Errorf(" Error while getting block details : %v", err.Error())
			return err
		}
		transactionCount, err := cl.TransactionCount(context.Background(), block.Hash())
		if err != nil {
			logging.Errorf(" Error while getting transaction count : %v", err.Error())
			return err
		}
		totalTransactions = totalTransactions + int(transactionCount)
	}
	logging.Infof("Total number of transaction between %v and %v is %v.", startBlock, endBlock, totalTransactions)

	if totalTransactions >= maxTxLoad {
		logging.Infof("Transaction load is higher than the %v for %v blocks, Please check the blockchain.", maxTxLoad, blockDifferenceForMaxTxLoad)
		// TODO: send email
		emaiMessage := "Alert, Blockchain has reached its maximum threshold ! \n\n" +
			"Maximum threshold per " + config.BlockDifferenceForMaxTxLoad + " blocks is " +
			config.MaxTxLoad + "\n" + "Number of transactions between " + strconv.FormatInt(startBlock, 10) +
			" and " + strconv.FormatInt(endBlock, 10) + " was " + strconv.Itoa(totalTransactions)
		err := email.SendEmail(constants.EmailSubject, emaiMessage)
		if err != nil {
			logging.Error("Error while sending Email")
			return err
		}
	}

	logging.Info("CheckForMaxTxLoad end")
	return nil
}
