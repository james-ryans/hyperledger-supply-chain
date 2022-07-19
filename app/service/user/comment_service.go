package userservice

import (
	"github.com/google/uuid"
	usermodel "github.com/meneketehe/hehe/app/model/user"
)

type commentService struct {
	CommentRepository usermodel.CommentRepository
}

type CommentServiceConfig struct {
	CommentRepository usermodel.CommentRepository
}

func NewCommentService(c *CommentServiceConfig) usermodel.CommentService {
	return &commentService{
		CommentRepository: c.CommentRepository,
	}
}

func (s *commentService) GetAllCommentByRiceID(riceID string) ([]*usermodel.Comment, error) {
	return s.CommentRepository.FindAllByRiceID(riceID)
}

func (s *commentService) WriteComment(comment *usermodel.Comment) (*usermodel.Comment, error) {
	comment.ID = uuid.New().String()

	if err := s.CommentRepository.Create(comment); err != nil {
		return nil, err
	}

	return comment, nil
}
