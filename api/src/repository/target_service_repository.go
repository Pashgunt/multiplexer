package repository

import (
	"context"
	"database/sql"
	"errors"
	"transport/api/src/model"
	"transport/api/src/vo"
	"transport/internal/infrastructure/db"
)

type ITargetServiceRepository interface {
	Save(ctx context.Context, targetService *model.TargetService) error
	CheckIssetServiceName(ctx context.Context, serviceName vo.ServiceName) (string, error)
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
