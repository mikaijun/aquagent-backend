package repositoryimpl

import (
	"context"

	"github.com/mikaijun/anli/pkg/infrastructure"

	"github.com/mikaijun/anli/pkg/domain/model"
	"github.com/mikaijun/anli/pkg/domain/repository"
)

type questionRepositoryImpl struct {
	db infrastructure.DBTX
}

func NewQuestionRepositoryImpl(db infrastructure.DBTX) repository.QuestionRepository {
	return &questionRepositoryImpl{db: db}
}

func (ri *questionRepositoryImpl) CreateQuestion(ctx context.Context, question *model.Question) (*model.Question, error) {
	var lastInsertId int
	query := "INSERT INTO questions (user_id, title, content, file_path, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7) returning id"
	err := ri.db.QueryRowContext(
		ctx,
		query,
		question.UserID,
		question.Title,
		question.Content,
		question.FilePath,
		question.CreatedAt,
		question.UpdatedAt,
		nil,
	).Scan(&lastInsertId)
	if err != nil {
		return &model.Question{}, err
	}

	question.ID = int64(lastInsertId)
	return question, nil
}

func (ri *questionRepositoryImpl) GetQuestions(ctx context.Context, userId int64) ([]*model.Question, error) {
	var questions []*model.Question
	query := "SELECT id, user_id, title, content, file_path, created_at, updated_at FROM questions WHERE deleted_at IS NULL AND user_id = $1 ORDER BY created_at DESC"
	rows, err := ri.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		question := &model.Question{}
		err := rows.Scan(
			&question.ID,
			&question.UserID,
			&question.Title,
			&question.Content,
			&question.FilePath,
			&question.CreatedAt,
			&question.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}
	return questions, nil
}
