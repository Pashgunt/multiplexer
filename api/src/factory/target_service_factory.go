package factory

import (
	"time"
	"transport/api/src/command"
	apidtodb "transport/api/src/dto/db"
	"transport/api/src/model"
	"transport/api/src/vo"

	"github.com/google/uuid"
)

type ITargetServiceFactory interface {
	Create(command command.CreateTargetServiceCommand) (*model.TargetService, error)
	CreateFromDb(dto apidtodb.TargetServiceDbDto) (*model.TargetService, error)
}

type TargetServiceFactory struct {
}

func NewTargetServiceFactory() ITargetServiceFactory {
	return TargetServiceFactory{}
}

func (f TargetServiceFactory) Create(command command.CreateTargetServiceCommand) (*model.TargetService, error) {
	serviceName, err := vo.NewServiceName(command.Dto.ServiceName, vo.ServiceNameValidateLevelStrict)

	if err != nil {
		return nil, err
	}

	baseURL, err := vo.NewBaseURL(command.Dto.BaseURL, vo.BaseURLValidateLevelStrict)

	if err != nil {
		return nil, err
	}

	date := time.Now()

	return model.NewTargetService(
		uuid.New(),
		serviceName,
		command.Dto.Description,
		baseURL,
		command.Dto.IsActive,
		&date,
		&date,
	), nil
}

func (f TargetServiceFactory) CreateFromDb(dto apidtodb.TargetServiceDbDto) (*model.TargetService, error) {
	serviceName, err := vo.NewServiceName(dto.ServiceName, vo.ServiceNameValidateLevelNone)

	if err != nil {
		return nil, err
	}

	baseURL, err := vo.NewBaseURL(dto.BaseURL, vo.BaseURLValidateLevelNone)

	if err != nil {
		return nil, err
	}

	date := time.Now()

	return model.NewTargetService(
		uuid.MustParse(dto.ID),
		serviceName,
		dto.Description,
		baseURL,
		dto.IsActive,
		nil,
		&date,
	), nil
}
