package repository

import (
	"context"

	"github.com/mikaijun/aquagent/pkg/domain/model"
)

type QuestionRepository interface {
	CreateQuestion(ctx context.Context, question *model.Question) (*model.Question, error)
	GetQuestions(ctx context.Context, userId int64) ([]*model.Question, error)
	GetQuestion(ctx context.Context, questionId int64) (*model.Question, error)
	UpdateQuestion(ctx context.Context, question *model.Question) (*model.Question, error)
}
