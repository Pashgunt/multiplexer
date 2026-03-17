package factory

import (
	"time"
	"transport/api/src/command"
	"transport/api/src/model"
	"transport/api/src/vo"

	"github.com/google/uuid"
)

type ITargetServiceFactory interface {
	Create(command command.CreateTargetServiceCommand) (*model.TargetService, error)
}

type TargetServiceFactory struct {
}

func NewTargetServiceFactory() ITargetServiceFactory {
	return TargetServiceFactory{}
}

func (f TargetServiceFactory) Create(command command.CreateTargetServiceCommand) (*model.TargetService, error) {
	serviceName, err := vo.NewServiceName(command.Dto.ServiceName)

	if err != nil {
		return nil, err
	}

	baseUrl, err := vo.NewBaseUrl(command.Dto.BaseUrl)

	if err != nil {
		return nil, err
	}

	return model.NewTargetService(
		uuid.New(),
		serviceName,
		command.Dto.Description,
		baseUrl,
		command.Dto.IsActive,
		time.Now(),
		time.Now(),
	), nil
}
