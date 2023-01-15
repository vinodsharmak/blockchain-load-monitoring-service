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

	// log.Info("Blockchain monitoring start")
	log.Info("Blockchain monitoring start")

	s := service.Service{}

	err = s.Configure()
	if err != nil {
		log.Error("Error while configuring the sub-services : ", err)
	}
	for {
		err = s.StartTxCountMonitoring()
		if err != nil {
			log.Errorf("Error while blockchain monitoring", err.Error())
		}
		time.Sleep(time.Second * 10)
	}

	// result, err := helpers.TxPoolContent()
	// if err != nil {
	// 	log.Errorf("error in txpool Content: %w", err)
	// }

	// bytes, err := json.MarshalIndent(result, " ", " ")
	// if err != nil {
	// 	log.Errorf("error in marshal indent: %w", err)
	// }
	// fmt.Println("RESPONSE: ", string(bytes))

	err = s.StartPendingAndQueuedTxMonitoring()
	if err != nil {
		log.Error("Error while pending and queued tx monitoring : ", err)
	}
	time.Sleep(10000)
}
