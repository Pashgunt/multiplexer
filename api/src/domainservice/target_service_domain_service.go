package domainservice

import (
	"context"
	apierror "transport/api/src/error"
	"transport/api/src/factory"
	"transport/api/src/model"
	"transport/api/src/repository"

	"github.com/google/uuid"
)

type ITargetDomainService interface {
	CanCreate(ctx context.Context, targetService *model.TargetService) error
	CanDelete(ctx context.Context, uuid uuid.UUID) (*model.TargetService, error)
}

type TargetDomainService struct {
	repository           repository.ITargetServiceRepository
	targetServiceFactory factory.ITargetServiceFactory
}

func NewTargetDomainService(
	repository repository.ITargetServiceRepository,
	targetServiceFactory factory.ITargetServiceFactory,
) ITargetDomainService {
	return &TargetDomainService{
		repository:           repository,
		targetServiceFactory: targetServiceFactory,
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

func (s *TargetDomainService) CanDelete(ctx context.Context, uuid uuid.UUID) (*model.TargetService, error) {
	targetServiceDto, err := s.repository.FindBy(ctx, uuid)

	if err != nil {
		return nil, err
	}

	if targetServiceDto == nil {
		return nil, apierror.NewNotFoundError("target_service")
	}

	return s.targetServiceFactory.CreateFromDb(*targetServiceDto)
}
