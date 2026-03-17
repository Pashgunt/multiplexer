package model

import (
	"time"
	"transport/api/src/vo"

	"github.com/google/uuid"
)

type TargetService struct {
	id          uuid.UUID
	serviceName vo.ServiceName
	description string
	baseUrl     vo.BaseUrl
	isActive    bool
	createdAt   time.Time
	updatedAt   time.Time
}

func (t *TargetService) Id() uuid.UUID {
	return t.id
}

func (t *TargetService) SetId(id uuid.UUID) {
	t.id = id
}

func (t *TargetService) ServiceName() vo.ServiceName {
	return t.serviceName
}

func (t *TargetService) SetServiceName(serviceName vo.ServiceName) {
	t.serviceName = serviceName
}

func (t *TargetService) Description() string {
	return t.description
}

func (t *TargetService) SetDescription(description string) {
	t.description = description
}

func (t *TargetService) BaseUrl() vo.BaseUrl {
	return t.baseUrl
}

func (t *TargetService) SetBaseUrl(baseUrl vo.BaseUrl) {
	t.baseUrl = baseUrl
}

func (t *TargetService) IsActive() bool {
	return t.isActive
}

func (t *TargetService) SetIsActive(isActive bool) {
	t.isActive = isActive
}

func (t *TargetService) CreatedAt() time.Time {
	return t.createdAt
}

func (t *TargetService) SetCreatedAt(createdAt time.Time) {
	t.createdAt = createdAt
}

func (t *TargetService) UpdatedAt() time.Time {
	return t.updatedAt
}

func (t *TargetService) SetUpdatedAt(updatedAt time.Time) {
	t.updatedAt = updatedAt
}

func NewTargetService(
	id uuid.UUID,
	serviceName vo.ServiceName,
	description string,
	baseUrl vo.BaseUrl,
	isActive bool,
	createdAt time.Time,
	updatedAt time.Time,
) *TargetService {
	return &TargetService{
		id:          id,
		serviceName: serviceName,
		description: description,
		baseUrl:     baseUrl,
		isActive:    isActive,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}
