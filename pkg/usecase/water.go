package usecase

import (
	"context"
	"time"

	"github.com/mikaijun/aquagent/pkg/domain/model"
	"github.com/mikaijun/aquagent/pkg/domain/repository"
)

type WaterUseCase interface {
	Get(c context.Context, id int64) (*model.Water, error)
	GetAll(c context.Context, userId int64) ([]*model.Water, error)
	Create(c context.Context, water *model.Water) (*model.Water, error)
	Update(c context.Context, water *model.Water) (*model.Water, error)
	Delete(c context.Context, id int64) error
}

type waterUseCase struct {
	repository repository.WaterRepository
	timeout    time.Duration
}

func NewWaterUseCase(waterRepo repository.WaterRepository) WaterUseCase {
	return &waterUseCase{
		repository: waterRepo,
		timeout:    time.Duration(2) * time.Second,
	}
}

func (uc *waterUseCase) Get(c context.Context, id int64) (*model.Water, error) {
	ctx, cancel := context.WithTimeout(c, uc.timeout)
	defer cancel()

	water, err := uc.repository.GetWater(ctx, id)
	if err != nil {
		return nil, err
	}

	return water, nil

}

func (uc *waterUseCase) GetAll(c context.Context, userId int64) ([]*model.Water, error) {
	ctx, cancel := context.WithTimeout(c, uc.timeout)
	defer cancel()

	waters, err := uc.repository.GetWaters(ctx, userId)
	if err != nil {
		return nil, err
	}

	return waters, nil
}

func (uc *waterUseCase) Create(c context.Context, water *model.Water) (*model.Water, error) {
	ctx, cancel := context.WithTimeout(c, uc.timeout)
	defer cancel()

	water, err := uc.repository.CreateWater(ctx, water)
	if err != nil {
		return nil, err
	}

	return water, nil
}

func (uc *waterUseCase) Update(c context.Context, water *model.Water) (*model.Water, error) {
	ctx, cancel := context.WithTimeout(c, uc.timeout)
	defer cancel()

	_, err := uc.repository.GetWater(ctx, water.ID)
	if err != nil {
		return nil, err
	}

	water, err = uc.repository.UpdateWater(ctx, water)
	if err != nil {
		return nil, err
	}

	return water, nil
}

func (uc *waterUseCase) Delete(c context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(c, uc.timeout)
	defer cancel()

	_, err := uc.repository.GetWater(ctx, id)
	if err != nil {
		return err
	}

	err = uc.repository.DeleteWater(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
