package model

type Question struct {
	ID        int64
	UserID    int64
	Title     string
	Content   string
	FilePath  string
	CreatedAt string
	UpdatedAt string
	DeletedAt string
}
