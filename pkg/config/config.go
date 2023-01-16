package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/antigloss/go/logger"
	"github.com/joho/godotenv"
)

var (
	BlockchainURL                string
	Logger                       *logger.Logger
	MaxTxLoad                    string
	BlockDifferenceForMaxTxLoad  string
	MaxTxPerBlock                string
	MaxTxPending                 string
	ChainID                      string
	MaxGasUsedPerBlock           string
	BlockDifferenceForMaxGasUsed string
	TimeIntervalForSubService    int
	TxpoolTimeLimit              int
)

// ReadConfig reads config file into the Config struct and returns it
func ReadConfig() error {
	var err error
	Logger, err = logger.New(&logger.Config{
		LogDir:          "./logs",
		LogFileMaxSize:  50,
		LogFileMaxNum:   10,
		LogFileNumToDel: 3,
		LogLevel:        logger.LogLevelInfo,
		LogDest:         logger.LogDestFile,
		Flag:            logger.ControlFlagLogLineNum,
	})

	if err != nil {
		return errors.New("unable to intialize logger")
	}
	err = godotenv.Load(".env")
	if err != nil {
		return errors.New("unable to load .env file")
	}

	blockchainURL, exists := os.LookupEnv("BLOCKCHAIN_URL")
	if !exists || blockchainURL == "" {
		return errors.New("blockchian URL cannot be empty")
	}
	BlockchainURL = blockchainURL
	Logger.Infof("BlockchainURL: %v", BlockchainURL)

	maxTxLoad, exists := os.LookupEnv("MAX_TX_LOAD")
	if !exists || maxTxLoad == "" {
		return errors.New("MAX_TX_LOAD cannot be empty")
	}
	MaxTxLoad = maxTxLoad
	Logger.Infof("MaxTxLoad: %v", MaxTxLoad)

	blockDifferenceForMaxTxLoad, exists := os.LookupEnv("BLOCK_DIFFERENCE_FOR_MAX_TX_LOAD")
	if !exists || blockDifferenceForMaxTxLoad == "" {
		return errors.New("BLOCK_DIFFERENCE_FOR_MAX_TX_LOAD cannot be empty")
	}
	BlockDifferenceForMaxTxLoad = blockDifferenceForMaxTxLoad
	Logger.Infof("BlockDifferenceForMaxTxLoad: %v", BlockDifferenceForMaxTxLoad)

	maxTxPerBlock, exists := os.LookupEnv("MAX_TX_PER_BLOCK")
	if !exists || maxTxPerBlock == "" {
		return errors.New("MAX_TX_PER_BLOCK cannot be empty")
	}
	MaxTxPerBlock = maxTxPerBlock
	Logger.Infof("MaxTxPerBlock: %v", MaxTxPerBlock)

	maxTxPending, exists := os.LookupEnv("MAX_TX_PENDING")
	if !exists || maxTxPending == "" {
		return errors.New("MAX_TX_PENDING cannot be empty")
	}
	MaxTxPending = maxTxPending
	Logger.Infof("MaxTxPending: %v", MaxTxPending)

	chainID, exists := os.LookupEnv("CHAIN_ID")
	if !exists || chainID == "" {
		return errors.New("CHAIN_ID cannot be empty")
	}
	ChainID = chainID
	Logger.Infof("ChainID: %v", ChainID)

	maxGasUsedPerBlock, exists := os.LookupEnv("MAX_GAS_USED_PER_BLOCK")
	if !exists || maxGasUsedPerBlock == "" {
		return errors.New("MAX_GAS_USED_PER_BLOCK cannot be empty")
	}
	MaxGasUsedPerBlock = maxGasUsedPerBlock
	Logger.Infof("MaxGasUsedPerBlock: %v", MaxGasUsedPerBlock)

	blockDifferenceForMaxGasUsed, exists := os.LookupEnv("BLOCK_DIFFERENCE_FOR_MAX_GAS_USED")
	if !exists || blockDifferenceForMaxGasUsed == "" {
		return errors.New("BLOCK_DIFFERENCE_FOR_MAX_GAS_USED cannot be empty")
	}
	BlockDifferenceForMaxGasUsed = blockDifferenceForMaxGasUsed
	Logger.Infof("BlockDifferenceForMaxGasUsed: %v", BlockDifferenceForMaxGasUsed)

	timeIntervalForSubService, exists := os.LookupEnv("TIME_INTERVAL_FOR_SUB_SERVICES")
	if !exists || timeIntervalForSubService == "" {
		return errors.New("TIME_INTERVAL_FOR_SUB_SERVICES cannot be empty")
	}
	TimeIntervalForSubService, err = strconv.Atoi(timeIntervalForSubService)
	if err != nil {
		return errors.New("unable to parse timeIntervalForSubService from string to integer, invalid format")
	}
	Logger.Infof("TimeIntervalForSubService: %v", TimeIntervalForSubService)

	txpoolTimeLimit, exists := os.LookupEnv("TXPOOL_TIME_LIMIT_IN_SECONDS")
	if !exists || chainID == "" {
		return errors.New("TXPOOL_TIME_LIMIT_IN_SECONDS cannot be empty")
	}
	TxpoolTimeLimit, err = strconv.Atoi(txpoolTimeLimit)
	if err != nil {
		return errors.New("unable to parse txpoolTimeLimit from string to integer, invalid format")
	}
	Logger.Infof("TxpoolTimeLimit: %v", TxpoolTimeLimit)

	return nil

}
