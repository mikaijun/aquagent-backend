package repositoryimpl

import (
	"context"
	"errors"
	"time"

	"github.com/mikaijun/aquagent/pkg/infrastructure"

	"github.com/mikaijun/aquagent/pkg/domain/model"
	"github.com/mikaijun/aquagent/pkg/domain/repository"
)

type waterRepositoryImpl struct {
	db infrastructure.DBTX
}

var createdAt time.Time
var updatedAt time.Time

func NewWaterRepositoryImpl(db infrastructure.DBTX) repository.WaterRepository {
	return &waterRepositoryImpl{db: db}
}

func (ri *waterRepositoryImpl) CreateWater(ctx context.Context, water *model.Water) (*model.Water, error) {
	var lastInsertId int
	query := "INSERT INTO waters (user_id, volume, created_at, updated_at) VALUES ($1, $2, $3, $4) returning id"
	err := ri.db.QueryRowContext(
		ctx,
		query,
		water.UserID,
		water.Volume,
		water.CreatedAt,
		water.UpdatedAt,
	).Scan(&lastInsertId)
	if err != nil {
		return &model.Water{}, err
	}

	water.ID = int64(lastInsertId)
	return water, nil
}

func (ri *waterRepositoryImpl) GetWaters(ctx context.Context, userId int64, filter map[string]interface{}) ([]*model.Water, error) {
	var waters []*model.Water = []*model.Water{}
	query := "SELECT id, user_id, volume, created_at, updated_at FROM waters WHERE user_id = $1"
	args := []interface{}{userId}

	// 日程指定 (2024/01/01 など)
	if date, ok := filter["date"].(string); ok {
		query += " AND DATE(created_at) = $2"
		args = append(args, date)
	}

	// 期間指定 (2024年の5月 など)
	if month, ok := filter["month"].(string); ok {
		query += " AND DATE_TRUNC('month', created_at) = DATE_TRUNC('month', TO_DATE($2, 'YYYY-MM'))"
		args = append(args, month)
	}

	query += " ORDER BY created_at DESC"

	rows, err := ri.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		water := &model.Water{}
		err := rows.Scan(
			&water.ID,
			&water.UserID,
			&water.Volume,
			&createdAt,
			&updatedAt,
		)
		water.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		water.UpdatedAt = updatedAt.Format("2006-01-02 15:04:05")
		if err != nil {
			return nil, err
		}
		waters = append(waters, water)
	}
	return waters, nil
}

func (ri *waterRepositoryImpl) GetWater(ctx context.Context, waterId int64) (*model.Water, error) {
	water := &model.Water{}
	query := "SELECT id, user_id, volume, created_at, updated_at FROM waters WHERE id = $1"

	err := ri.db.QueryRowContext(ctx, query, waterId).Scan(
		&water.ID,
		&water.UserID,
		&water.Volume,
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
	query := "UPDATE waters SET volume = $1, updated_at = $2 WHERE id = $3"
	_, err := ri.db.ExecContext(
		ctx,
		query,
		water.Volume,
		water.UpdatedAt,
		water.ID,
	)
	if err != nil {
		return &model.Water{}, err
	}
	return water, nil
}

func (ri *waterRepositoryImpl) DeleteWater(ctx context.Context, waterId int64) error {
	query := "DELETE FROM waters WHERE id = $1"
	_, err := ri.db.ExecContext(ctx, query, waterId)
	if err != nil {
		return err
	}
	return nil
}
