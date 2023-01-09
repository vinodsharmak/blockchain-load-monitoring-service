package service

func StartMonitoring() error {
	err := checkTxLoad()
	if err != nil {
		return err
	}
	return nil
}
