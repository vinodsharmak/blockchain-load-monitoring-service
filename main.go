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
		logrus.Errorln("Error while configuring blockchian load monitoring service : ", err)
		return
	}

	err = email.Config()
	if err != nil {
		logrus.Errorf("unable to configure ses: ", err)
		return
	}

	logging := config.Logger

	logging.Info("Blockchain monitoring start")
	err = service.StartMonitoring()
	if err != nil {
		logging.Errorf("Error while blockchain monitoring", err.Error())
		logging.Info("Blockchain monitoring end")
		return
	}
	logging.Info("Blockchain monitoring end")
}
