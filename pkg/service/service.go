package service

import "bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/config"

func StartMonitoring() error {
	logging := config.Logger

	logging.Info("Checking for max TX load...")
	err := CheckForMaxTxLoad()
	if err != nil {
		logging.Errorf("Error while checking for maximum load.")
		return err
	}

	return nil
}
