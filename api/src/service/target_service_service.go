package service

import (
	"context"
	"encoding/json"
	"time"
	"transport/api/src/command"
	"transport/api/src/domainservice"
	apidtodb "transport/api/src/dto/db"
	"transport/api/src/factory"
	"transport/api/src/repository"

	"github.com/redis/go-redis/v9"
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
	redis         *redis.Client
}

func NewTargetServiceService(
	repository repository.ITargetServiceRepository,
	targetServiceFactory factory.ITargetServiceFactory,
	redis *redis.Client,
	domainService domainservice.ITargetDomainService,
) ITargetServiceService {
	return &TargetServiceService{
		repository:    repository,
		factory:       targetServiceFactory,
		domainService: domainService,
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
	value, err := s.redis.Get(ctx, command.ID.String()).Bytes()

	if err == nil {
		var dto apidtodb.TargetServiceDbDto

		if err = json.Unmarshal(value, &dto); err != nil {
			return &dto, err
		}

		return &dto, nil
	}

	result, err := s.repository.FindBy(ctx, command.ID)

	s.redis.Set(ctx, command.ID.String(), result, 5*time.Minute)

	return result, err
}
