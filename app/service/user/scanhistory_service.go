package userservice

import usermodel "github.com/meneketehe/hehe/app/model/user"

type scanHistoryService struct {
	ScanHistoryRepository usermodel.ScanHistoryRepository
}

type ScanHistoryServiceConfig struct {
	ScanHistoryRepository usermodel.ScanHistoryRepository
}

func NewScanHistoryService(c *ScanHistoryServiceConfig) usermodel.ScanHistoryService {
	return &scanHistoryService{
		ScanHistoryRepository: c.ScanHistoryRepository,
	}
}

func (s *scanHistoryService) GetAllScanHistories(userID string) ([]*usermodel.ScanHistory, error) {
	return s.ScanHistoryRepository.FindAll(userID)
}
