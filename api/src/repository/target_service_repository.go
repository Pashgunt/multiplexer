package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	apidtodb "transport/api/src/dto/db"
	"transport/api/src/model"
	"transport/api/src/vo"

	"github.com/google/uuid"
)

type ITargetServiceRepository interface {
	Save(ctx context.Context, targetService *model.TargetService) error
	CheckIssetServiceName(ctx context.Context, serviceName vo.ServiceName) (string, error)
	Delete(ctx context.Context, uuid uuid.UUID) error
	FindBy(ctx context.Context, uuid uuid.UUID) (*apidtodb.TargetServiceDbDto, error)
}

const (
	targetServiceTableName = "target_services"
)

type TargetServiceRepository struct {
	connection *sql.DB
}

func NewTargetServiceRepository(connection *sql.DB) ITargetServiceRepository {
	return &TargetServiceRepository{connection: connection}
}

func (r TargetServiceRepository) CheckIssetServiceName(
	ctx context.Context,
	serviceName vo.ServiceName,
) (string, error) {
	var result string

	err := r.
		connection.
		QueryRowContext(
			ctx,
			fmt.Sprintf("select ts.id from %s ts where ts.service_name = $1", targetServiceTableName),
			serviceName.Value(),
		).
		Scan(&result)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		return "", err
	}

	return result, nil
}

func (r TargetServiceRepository) Save(ctx context.Context, targetService *model.TargetService) error {
	_, err := r.connection.ExecContext(
		ctx,
		fmt.Sprintf("insert into %s (id, service_name, description, base_url, is_active, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7)", targetServiceTableName),
		targetService.ID(),
		targetService.ServiceName().Value(),
		targetService.Description(),
		targetService.BaseURL().Value(),
		targetService.IsActive(),
		targetService.CreatedAt(),
		targetService.UpdatedAt(),
	)

	return err
}

func (r TargetServiceRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	_, err := r.connection.ExecContext(
		ctx,
		fmt.Sprintf("delete from %s where id = $1", targetServiceTableName),
		uuid,
	)

	return err
}

func (r TargetServiceRepository) FindBy(ctx context.Context, uuid uuid.UUID) (*apidtodb.TargetServiceDbDto, error) {
	row := r.connection.QueryRowContext(ctx,
		fmt.Sprintf("SELECT ts.id, ts.service_name, ts.description, ts.base_url, ts.is_active FROM %s ts WHERE ts.id = $1", targetServiceTableName),
		uuid.String())

	var targetServiceDbDto apidtodb.TargetServiceDbDto

	err := row.Scan(
		&targetServiceDbDto.ID,
		&targetServiceDbDto.ServiceName,
		&targetServiceDbDto.Description,
		&targetServiceDbDto.BaseURL,
		&targetServiceDbDto.IsActive,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &targetServiceDbDto, nil
}
