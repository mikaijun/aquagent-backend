package repository

import (
	"context"

	"github.com/mikaijun/anli/pkg/domain/model"
)

type QuestionRepository interface {
	CreateQuestion(ctx context.Context, question *model.Question) (*model.Question, error)
}