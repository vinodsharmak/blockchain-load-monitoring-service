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
	err = s.Configure()
	if err != nil {
		log.Errorf("error in configuring blocktx count service", err)
	}

	for {
		err = s.StartTxCountMonitoring()
		if err != nil {
			log.Errorf("Error while blockchain monitoring", err.Error())
		}
		time.Sleep(10000)
	}
}
