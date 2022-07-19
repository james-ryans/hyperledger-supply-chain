package usermodel

import "time"

type Comment struct {
	ID        string    `json:"id"`
	RiceID    string    `json:"rice_id"`
	UserName  string    `json:"user_name"`
	Text      string    `json:"text"`
	CommentAt time.Time `json:"comment_at"`
}

type CommentService interface {
	GetAllCommentByRiceID(riceID string) ([]*Comment, error)
	WriteComment(comment *Comment) (*Comment, error)
}

type CommentRepository interface {
	FindAllByRiceID(riceID string) ([]*Comment, error)
	Create(comment *Comment) error
}
