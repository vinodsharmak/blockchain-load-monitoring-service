package main

import (
	"time"

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

	s := service.Service{}
	s := service.GasUsedService{}
	
	err = s.Configure()
	if err != nil {
		log.Error("Error while configuring the sub-services : ", err)
	}

	err = s.StartPendingAndQueuedTxMonitoring()
	if err != nil {
		log.Error("Error while pending and queued tx monitoring : ", err)

	err = s.StartGasUsedtMonitoring()
	if err != nil {
		log.Errorf("Error while gas used monitoring", err.Error())
	}
	
}
