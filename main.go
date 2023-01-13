package main

import (
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/config"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/service"
	"github.com/sirupsen/logrus"
)

func main() {
	err := config.ReadConfig()
	if err != nil {
		logrus.Errorln("Error while configuring blockchian load monitoring service : ", err)
		return
	}

	log := config.Logger

	log.Info("Blockchain monitoring start")

	err = service.StartPendingAndQueuedTxMonitoring()
	if err != nil {
		log.Errorf("Error while pending and queued tx monitoring : ", err)
	}
}
