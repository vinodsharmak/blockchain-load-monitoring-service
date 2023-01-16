package main

import (
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/config"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/email"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/service"
	"github.com/sirupsen/logrus"
)

func main() {
	err := config.ReadConfig()
	if err != nil {
		logrus.Errorln("read config: ", err)
		return
	}

	log := config.Logger

	err = email.Config()
	if err != nil {
		log.Errorf("unable to configure ses: %v", err)
		return
	}

	s := service.Service{}
	err = s.Configure()
	if err != nil {
		log.Errorf("configure %v: ", err)
	}

	err = s.BlockchainMonitoringService()
	if err != nil {
		log.Errorf("Blockchain Monitoring Service: %s ", err.Error())
	}
}
