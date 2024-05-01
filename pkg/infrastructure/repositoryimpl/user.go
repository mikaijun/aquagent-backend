package repositoryimpl

import (
	"context"

	"github.com/mikaijun/anli/pkg/infrastructure"

	"github.com/mikaijun/anli/pkg/domain/model"
	"github.com/mikaijun/anli/pkg/domain/repository"
)

type repositoryImpl struct {
	db infrastructure.DBTX
}

func NewUserRepositoryImpl(db infrastructure.DBTX) repository.UserRepository {
	return &repositoryImpl{db: db}
}

func (ri *repositoryImpl) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	var lastInsertId int
	query := "INSERT INTO users(username, email, password) VALUES ($1, $2, $3) returning id"
	err := ri.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password).Scan(&lastInsertId)
	if err != nil {
		return &model.User{}, err
	}

	user.ID = int64(lastInsertId)
	return user, nil
}

func (ri *repositoryImpl) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	u := model.User{}
	query := "SELECT id, username, email, password FROM users WHERE email = $1"
	err := ri.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Username, &u.Email, &u.Password)
	if err != nil {
		return &model.User{}, nil
	}

	return &u, nil
}

func (ri *repositoryImpl) GetUserById(ctx context.Context, id int64) (*model.User, error) {
	u := model.User{}
	query := "SELECT id, username, email, password FROM users WHERE id = $1"
	err := ri.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Username, &u.Email, &u.Password)
	if err != nil {
		return &model.User{}, nil
	}

	return &u, nil
}
