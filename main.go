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
	s := service.Service{}
	err = s.Configure()
	if err != nil {
		log.Error("configure: ", err)
	}

	err = s.StartBlockchainMonitoringService()
	if err != nil {
		log.Errorf("Error while blockchain monitoring", err.Error())
	}
}
