package repository

import (
	"context"

	"github.com/mikaijun/anli/pkg/domain/model"
)

type QuestionRepository interface {
	CreateQuestion(ctx context.Context, question *model.Question) (*model.Question, error)
	GetQuestions(ctx context.Context, userId int64) ([]*model.Question, error)
}
