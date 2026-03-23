package service

import (
	"context"
	"time"
	"transport/api/src/command"
	"transport/api/src/domainservice"
	apidtodb "transport/api/src/dto/db"
	"transport/api/src/factory"
	"transport/api/src/repository"
	"transport/internal/infrastructure/redis"
)

type ITargetServiceService interface {
	Create(ctx context.Context, command command.CreateTargetServiceCommand) (string, error)
	Delete(ctx context.Context, command command.DeleteTargetServiceCommand) error
	Get(ctx context.Context, command command.GetTargetServiceCommand) (*apidtodb.TargetServiceDbDto, error)
}

type TargetServiceService struct {
	repository    repository.ITargetServiceRepository
	factory       factory.ITargetServiceFactory
	domainService domainservice.ITargetDomainService
	redis         redis.IRedis
}

func NewTargetServiceService(
	repository repository.ITargetServiceRepository,
	factory factory.ITargetServiceFactory,
	redis redis.IRedis,
) *TargetServiceService {
	return &TargetServiceService{
		repository:    repository,
		factory:       factory,
		domainService: domainservice.NewTargetDomainService(repository, factory),
		redis:         redis,
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

	return entity.ID().String(), s.repository.Save(ctx, entity)
}

func (s TargetServiceService) Delete(ctx context.Context, command command.DeleteTargetServiceCommand) error {
	targetService, err := s.domainService.CanDelete(ctx, command.ID)

	if err != nil {
		return err
	}

	return s.repository.Delete(ctx, targetService.ID())
}

func (s TargetServiceService) Get(ctx context.Context, command command.GetTargetServiceCommand) (*apidtodb.TargetServiceDbDto, error) {
	var dto apidtodb.TargetServiceDbDto
	err := s.redis.Get(ctx, command.ID.String(), dto)

	if err == nil {
		return &dto, nil
	}

	result, err := s.repository.FindBy(ctx, command.ID)

	s.redis.Set(ctx, command.ID.String(), result, 5*time.Minute)

	return result, err
}
