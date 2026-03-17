package service

import (
	"context"
	"transport/api/src/command"
	"transport/api/src/domain_service"
	"transport/api/src/factory"
	"transport/api/src/repository"
)

type ITargetServiceService interface {
	Create(ctx context.Context, command command.CreateTargetServiceCommand) (string, error)
}

type TargetServiceService struct {
	repository    repository.ITargetServiceRepository
	factory       factory.ITargetServiceFactory
	domainService domain_service.ITargetDomainService
}

func NewTargetServiceService(
	repository repository.ITargetServiceRepository,
	factory factory.ITargetServiceFactory,
) *TargetServiceService {
	return &TargetServiceService{
		repository:    repository,
		factory:       factory,
		domainService: domain_service.NewTargetDomainService(repository),
	}
}

func (s TargetServiceService) Create(ctx context.Context, command command.CreateTargetServiceCommand) (string, error) {
	entity, err := s.factory.Create(command)

	if err != nil {
		return "", err
	}

	if err = s.domainService.CanCreate(ctx, entity); err != nil {
		return "", err
	}

	return entity.Id().String(), s.repository.Save(ctx, entity)
}
