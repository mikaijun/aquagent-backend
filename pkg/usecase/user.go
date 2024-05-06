package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/mikaijun/anli/pkg/domain/model"
	"github.com/mikaijun/anli/pkg/domain/repository"
	"github.com/mikaijun/anli/pkg/util"
)

type UserUseCase interface {
	Signup(c context.Context, username, email, password string) (*model.User, error)
	Login(c context.Context, email, password string) (string, *model.User, error)
	Fetch(c context.Context, userId int64) (*model.User, error)
}

type userUseCase struct {
	repository repository.UserRepository
	timeout    time.Duration
}

func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return &userUseCase{
		repository: userRepo,
		timeout:    time.Duration(2) * time.Second,
	}
}

func (uc *userUseCase) Signup(c context.Context, username, email, password string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(c, uc.timeout)
	defer cancel()

	exsitUser, err := uc.repository.GetUserByEmail(ctx, email)

	if err != nil {
		return nil, &util.InternalServerError{Err: err}
	}

	if exsitUser.ID != 0 {
		return nil, &util.BadRequestError{Err: errors.New("user already exists")}
	}

	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return nil, &util.InternalServerError{Err: err}
	}

	u := &model.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	user, err := uc.repository.CreateUser(ctx, u)
	if err != nil {
		return nil, &util.InternalServerError{Err: err}
	}

	return user, nil
}

func (uc *userUseCase) Login(c context.Context, email, password string) (string, *model.User, error) {
	ctx, cancel := context.WithTimeout(c, uc.timeout)
	defer cancel()

	user, err := uc.repository.GetUserByEmail(ctx, email)
	if err != nil {
		return "", nil, &util.InternalServerError{Err: err}
	}
	if user.ID == 0 {
		return "", nil, &util.BadRequestError{Err: errors.New("user is not exist")}
	}

	err = util.CheckPassword(user.Password, password)
	if err != nil {
		return "", nil, &util.BadRequestError{Err: errors.New("password is incorrect")}
	}

	signedString, err := util.GenerateSignedString(user.ID, user.Username)
	if err != nil {
		return "", nil, &util.InternalServerError{Err: err}
	}

	return signedString, user, nil
}

func (uc *userUseCase) Fetch(c context.Context, userId int64) (*model.User, error) {
	ctx, cancel := context.WithTimeout(c, uc.timeout)
	defer cancel()

	user, err := uc.repository.GetUserById(ctx, userId)
	if err != nil {
		return nil, &util.InternalServerError{Err: err}
	}
	if user.ID == 0 {
		return nil, &util.BadRequestError{Err: errors.New("user is not exist")}
	}
	return user, nil
}
