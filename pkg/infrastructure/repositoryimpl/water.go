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

var drankAt time.Time

func NewWaterRepositoryImpl(db infrastructure.DBTX) repository.WaterRepository {
	return &waterRepositoryImpl{db: db}
}

func (ri *waterRepositoryImpl) CreateWater(ctx context.Context, water *model.Water) (*model.Water, error) {
	var lastInsertId int
	query := "INSERT INTO waters (user_id, volume, drank_at) VALUES ($1, $2, $3) returning id"
	err := ri.db.QueryRowContext(
		ctx,
		query,
		water.UserID,
		water.Volume,
		water.DrankAt,
	).Scan(&lastInsertId)
	if err != nil {
		return &model.Water{}, err
	}

	water.ID = int64(lastInsertId)
	return water, nil
}

func (ri *waterRepositoryImpl) GetWaters(ctx context.Context, userId int64, filter map[string]interface{}) ([]*model.Water, error) {
	var waters []*model.Water = []*model.Water{}
	query := "SELECT id, user_id, volume, drank_at FROM waters WHERE user_id = $1"
	args := []interface{}{userId}

	// 指定した日付以降(その日も含む)
	if start, ok := filter["start"].(string); ok {
		query += " AND DATE(drank_at) >= $2"
		args = append(args, start)
	}

	// 指定した日付以前(その日も含む)。startと一緒に使う想定
	if end, ok := filter["end"].(string); ok {
		query += " AND DATE(drank_at) <= $3"
		args = append(args, end)

	}

	query += " ORDER BY drank_at DESC"

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
			&drankAt,
		)
		water.DrankAt = drankAt.Format("2006-01-02 15:04:05")
		if err != nil {
			return nil, err
		}
		waters = append(waters, water)
	}
	return waters, nil
}

func (ri *waterRepositoryImpl) GetWater(ctx context.Context, waterId int64) (*model.Water, error) {
	water := &model.Water{}
	query := "SELECT id, user_id, volume, drank_at FROM waters WHERE id = $1"

	err := ri.db.QueryRowContext(ctx, query, waterId).Scan(
		&water.ID,
		&water.UserID,
		&water.Volume,
		&water.DrankAt,
	)

	if water.ID == 0 {
		return &model.Water{}, errors.New("water not found")
	}

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
