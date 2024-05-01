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
	query := `
		INSERT INTO questions (
			user_id,
			title,
			content,
			file_path,
			created_at,
			updated_at,
			deleted_at
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
	returning id
	`
	err := ri.db.QueryRowContext(
		ctx,
		query,
		question.UserID,
		question.Title,
		question.FilePath,
		question.CreatedAt,
		question.UpdatedAt,
		question.DeletedAt,
	).Scan(&lastInsertId)
	if err != nil {
		return &model.Question{}, err
	}

	question.ID = int64(lastInsertId)
	return question, nil
}
