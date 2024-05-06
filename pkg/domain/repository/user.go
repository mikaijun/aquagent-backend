package repository

import (
	"context"

	"github.com/mikaijun/aquagent/pkg/domain/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserById(ctx context.Context, id int64) (*model.User, error)
}
