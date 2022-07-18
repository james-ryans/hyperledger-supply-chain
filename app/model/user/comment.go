package usermodel

import "time"

type Comment struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	CommentAt time.Time `json:"comment_at"`
}
