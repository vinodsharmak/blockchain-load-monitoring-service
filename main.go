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
	for {
		err = service.StartMonitoring()
		if err != nil {
			log.Errorf("Error while blockchain monitoring", err.Error())
		}
		time.Sleep(10000)
	}

	// log.Info("Blockchain monitoring end")
}
