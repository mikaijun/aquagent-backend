package usecase

import (
	"context"
	"time"

	"github.com/mikaijun/anli/pkg/domain/model"
	"github.com/mikaijun/anli/pkg/domain/repository"
)

type QuestionUseCase interface {
	Create(c context.Context, question *model.Question) (*model.Question, error)
	GetAll(c context.Context, userId int64) ([]*model.Question, error)
}

type questionUseCase struct {
	repository repository.QuestionRepository
	timeout    time.Duration
}

func NewQuestionUseCase(questionRepo repository.QuestionRepository) QuestionUseCase {
	return &questionUseCase{
		repository: questionRepo,
		timeout:    time.Duration(2) * time.Second,
	}
}

func (uc *questionUseCase) Create(c context.Context, question *model.Question) (*model.Question, error) {
	ctx, cancel := context.WithTimeout(c, uc.timeout)
	defer cancel()

	question, err := uc.repository.CreateQuestion(ctx, question)
	if err != nil {
		return nil, err
	}

	return question, nil
}

func (uc *questionUseCase) GetAll(c context.Context, userId int64) ([]*model.Question, error) {
	ctx, cancel := context.WithTimeout(c, uc.timeout)
	defer cancel()

	questions, err := uc.repository.GetQuestions(ctx, userId)
	if err != nil {
		return nil, err
	}

	return questions, nil
}
