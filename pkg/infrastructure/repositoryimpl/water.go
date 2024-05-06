package repositoryimpl

import (
	"context"
	"errors"

	"github.com/mikaijun/aquagent/pkg/infrastructure"

	"github.com/mikaijun/aquagent/pkg/domain/model"
	"github.com/mikaijun/aquagent/pkg/domain/repository"
)

type waterRepositoryImpl struct {
	db infrastructure.DBTX
}

func NewWaterRepositoryImpl(db infrastructure.DBTX) repository.WaterRepository {
	return &waterRepositoryImpl{db: db}
}

func (ri *waterRepositoryImpl) CreateWater(ctx context.Context, water *model.Water) (*model.Water, error) {
	var lastInsertId int
	query := "INSERT INTO waters (user_id, title, content, file_path, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7) returning id"
	err := ri.db.QueryRowContext(
		ctx,
		query,
		water.UserID,
		water.Title,
		water.Content,
		water.FilePath,
		water.CreatedAt,
		water.UpdatedAt,
		nil,
	).Scan(&lastInsertId)
	if err != nil {
		return &model.Water{}, err
	}

	water.ID = int64(lastInsertId)
	return water, nil
}

func (ri *waterRepositoryImpl) GetWaters(ctx context.Context, userId int64) ([]*model.Water, error) {
	var waters []*model.Water
	query := "SELECT id, user_id, title, content, file_path, created_at, updated_at FROM waters WHERE deleted_at IS NULL AND user_id = $1 ORDER BY created_at DESC"
	rows, err := ri.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		water := &model.Water{}
		err := rows.Scan(
			&water.ID,
			&water.UserID,
			&water.Title,
			&water.Content,
			&water.FilePath,
			&water.CreatedAt,
			&water.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		waters = append(waters, water)
	}
	return waters, nil
}

func (ri *waterRepositoryImpl) GetWater(ctx context.Context, waterId int64) (*model.Water, error) {
	water := &model.Water{}
	query := "SELECT id, user_id, title, content, file_path, created_at, updated_at FROM waters WHERE deleted_at IS NULL AND id = $1"

	err := ri.db.QueryRowContext(ctx, query, waterId).Scan(
		&water.ID,
		&water.UserID,
		&water.Title,
		&water.Content,
		&water.FilePath,
		&water.CreatedAt,
		&water.UpdatedAt,
	)

	if water.ID == 0 {
		return &model.Water{}, errors.New("water not found")
	}

	if err != nil {
		return &model.Water{}, err
	}

	return water, nil
}

func (ri *waterRepositoryImpl) UpdateWater(ctx context.Context, water *model.Water) (*model.Water, error) {
	query := "UPDATE waters SET title = $1, content = $2, file_path = $3, updated_at = $4 WHERE id = $5 AND deleted_at IS NULL"
	_, err := ri.db.ExecContext(
		ctx,
		query,
		water.Title,
		water.Content,
		water.FilePath,
		water.UpdatedAt,
		water.ID,
	)
	if err != nil {
		return &model.Water{}, err
	}
	return water, nil
}
