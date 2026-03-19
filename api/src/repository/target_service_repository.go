package repository

import (
	"context"
	"database/sql"
	"errors"
	apidtodb "transport/api/src/dto/db"
	"transport/api/src/model"
	"transport/api/src/vo"
	"transport/internal/infrastructure/db"

	"github.com/google/uuid"
)

type ITargetServiceRepository interface {
	Save(ctx context.Context, targetService *model.TargetService) error
	CheckIssetServiceName(ctx context.Context, serviceName vo.ServiceName) (string, error)
	Delete(ctx context.Context, uuid uuid.UUID) error
	FindBy(ctx context.Context, uuid uuid.UUID) (*apidtodb.TargetServiceDbDto, error)
}

type TargetServiceRepository struct {
	connection db.DBInterface
}

func NewTargetServiceRepository(connection db.DBInterface) *TargetServiceRepository {
	return &TargetServiceRepository{connection: connection}
}

func (r TargetServiceRepository) CheckIssetServiceName(
	ctx context.Context,
	serviceName vo.ServiceName,
) (string, error) {
	var result string

	err := r.
		connection.
		Db().
		QueryRowContext(
			ctx,
			"select ts.id from target_services ts where ts.service_name = $1",
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
	_, err := r.connection.Db().ExecContext(
		ctx,
		"insert into target_services (id, service_name, description, base_url, is_active, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7)",
		targetService.Id(),
		targetService.ServiceName().Value(),
		targetService.Description(),
		targetService.BaseUrl().Value(),
		targetService.IsActive(),
		targetService.CreatedAt(),
		targetService.UpdatedAt(),
	)

	return err
}

func (r TargetServiceRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	_, err := r.connection.Db().ExecContext(
		ctx,
		"delete from target_services where id = $1", //todo make table name to const
		uuid,
	)

	return err
}

func (r TargetServiceRepository) FindBy(ctx context.Context, uuid uuid.UUID) (*apidtodb.TargetServiceDbDto, error) {
	row := r.connection.Db().QueryRowContext(ctx,
		"SELECT ts.id, ts.service_name, ts.description, ts.base_url, ts.is_active FROM target_services ts WHERE ts.id = $1",
		uuid.String())

	var targetServiceDbDto apidtodb.TargetServiceDbDto

	err := row.Scan(
		&targetServiceDbDto.Id,
		&targetServiceDbDto.ServiceName,
		&targetServiceDbDto.Description,
		&targetServiceDbDto.BaseUrl,
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
