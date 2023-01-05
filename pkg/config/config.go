package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/antigloss/go/logger"
	"github.com/joho/godotenv"
)

var (
	BlockchainURL               string
	Logger                      *logger.Logger
	MaxTxLoad                   string
	BlockDifferenceForMaxTxLoad string
	SenderEmail                 string
	ReceiverEmails              string
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
		panic(err)
	}
	err = godotenv.Load(".env")
	if err != nil {
		return errors.New("unable to load .env file")
	}

	blockchainURL, exists := os.LookupEnv("BLOCKCHAIN_URL")
	if !exists || blockchainURL == "" {
		return fmt.Errorf("%s", "Blockchian URL cannot be empty")
	}
	BlockchainURL = blockchainURL
	Logger.Infof("BlockchainURL: %v", blockchainURL)

	maxTxLoad, exists := os.LookupEnv("MAX_TX_LOAD")
	if !exists || maxTxLoad == "" {
		return fmt.Errorf("%s", "MAX_TX_LOAD cannot be empty")
	}
	MaxTxLoad = maxTxLoad
	Logger.Infof("MaxTxLoad: %v", maxTxLoad)

	blockDifferenceForMaxTxLoad, exists := os.LookupEnv("BLOCK_DIFFERENCE_FOR_MAX_TX_LOAD")
	if !exists || blockDifferenceForMaxTxLoad == "" {
		return fmt.Errorf("%s", "BLOCK_DIFFERENCE_FOR_MAX_TX_LOAD cannot be empty")
	}
	BlockDifferenceForMaxTxLoad = blockDifferenceForMaxTxLoad
	Logger.Infof("BlockDifferenceForMaxTxLoad: %v", blockDifferenceForMaxTxLoad)

	SenderEmail, exists = os.LookupEnv("SENDER_EMAIL")
	if !exists || SenderEmail == "" {
		Logger.Error("Sender email cannot be empty")
	}
	Logger.Info("Sender for Email Service ", SenderEmail)

	ReceiverEmails, exists = os.LookupEnv("RECEIVER_EMAILS")
	if !exists || ReceiverEmails == "" {
		logger.Error("Receiver emails cannot be empty")
	}
	Logger.Info("Receivers for Email Service ", ReceiverEmails)

	return nil
}
