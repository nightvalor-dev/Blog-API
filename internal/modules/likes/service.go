package likes

import (
	"context"
	"errors"
)

type LikeService struct {
	repo LikeRepository
}

func NewLikeService(repo LikeRepository) *LikeService {
	return &LikeService{repo: repo}
}

func (s *LikeService) Like(ctx context.Context, req LikeRequest) error {
	exists, err := s.repo.Exists(req.BlogId, req.UserId)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("already liked")
	}
	return s.repo.Add(req.BlogId, req.UserId)
}

func (s *LikeService) Unlike(ctx context.Context, req LikeRequest) error {
	exists, err := s.repo.Exists(req.BlogId, req.UserId)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("like not found")
	}
	return s.repo.Remove(req.BlogId, req.UserId)
}

func (s *LikeService) GetCount(ctx context.Context, blogId int) (LikeResponse, error) {
	count, err := s.repo.Count(blogId)
	if err != nil {
		return LikeResponse{}, err
	}
	return LikeResponse{BlogId: blogId, Count: count}, nil
}
