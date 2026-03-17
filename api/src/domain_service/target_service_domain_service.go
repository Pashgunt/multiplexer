package domain_service

import (
	"context"
	apierror "transport/api/src/error"
	"transport/api/src/model"
	"transport/api/src/repository"
)

type ITargetDomainService interface {
	CanCreate(ctx context.Context, targetService *model.TargetService) error
}

type TargetDomainService struct {
	repository repository.ITargetServiceRepository
}

func NewTargetDomainService(repository repository.ITargetServiceRepository) ITargetDomainService {
	return &TargetDomainService{
		repository: repository,
	}
}

func (s *TargetDomainService) CanCreate(ctx context.Context, targetService *model.TargetService) error {
	id, err := s.repository.CheckIssetServiceName(ctx, targetService.ServiceName())

	if err != nil {
		return err
	}

	if id != "" {
		return apierror.NewServiceNameUniqueError(id)
	}

	return nil
}
