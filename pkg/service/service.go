package service

func (s *Service) StartTxCountMonitoring() error {
	err := s.checkTxLoad()
	if err != nil {
		return err
	}
	return nil
}

func StartPendingAndQueuedTxMonitoring() error {
	err := CheckPendingAndQueuedTxCount()
	if err != nil {
		return err
	}
	return nil
}
