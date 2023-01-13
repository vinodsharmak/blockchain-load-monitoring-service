package main

import (
	"strconv"
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
		log.Error("Error while configuring the sub-services : ", err)
	}

	timeInterval, err := strconv.Atoi(config.TimeIntervalForSubService)
	if err != nil {
		log.Error("Error while string to integer conversion of timeInterval : ", err)
	}
	for {
		err = s.StartTxCountMonitoring()
		if err != nil {
			log.Error("Error while pending and queued tx monitoring : ", err)
		}
		time.Sleep(time.Duration(timeInterval) * time.Second)
		err = s.StartPendingAndQueuedTxMonitoring()
		if err != nil {
			log.Error("Error while pending and queued tx monitoring : ", err)
		}
		time.Sleep(time.Duration(timeInterval) * time.Second)
		err = s.StartGasUsedtMonitoring()
		if err != nil {
			log.Errorf("Error while gas used monitoring", err.Error())
		}
		time.Sleep(time.Duration(timeInterval) * time.Second)
	}
}
