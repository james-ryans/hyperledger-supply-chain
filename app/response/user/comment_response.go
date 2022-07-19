package userresponse

import (
	"time"

	usermodel "github.com/meneketehe/hehe/app/model/user"
)

type commentsResponse struct {
	ID        string    `json:"id"`
	UserName  string    `json:"user_name"`
	Text      string    `json:"text"`
	CommentAt time.Time `json:"comment_at"`
}

func CommentsResponse(comments []*usermodel.Comment) []*commentsResponse {
	res := make([]*commentsResponse, 0)
	for _, comment := range comments {
		res = append(res, &commentsResponse{
			ID:        comment.ID,
			UserName:  comment.UserName,
			Text:      comment.Text,
			CommentAt: comment.CommentAt,
		})
	}

	return res
}
