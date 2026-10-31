package service

import "github.com/LiddleChild/liddle-vpc/backend/internal/resources"

type IngestionService struct{}

func NewIngestionService() IngestionService {
	return IngestionService{}
}

func (svc IngestionService) GetModelIngestionInstruction() string {
	return resources.IngestionInstruction
}
