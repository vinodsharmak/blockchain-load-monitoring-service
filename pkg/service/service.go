package service

func (s *Service) StartTxCountMonitoring() error {
	err := s.checkTxLoad()
	if err != nil {
		return err
	}
	return nil
}
