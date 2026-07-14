package service

import (
    "errors"
    "fmt"
    "github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
    "github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/business/repository"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type PostService struct {
    repo *repository.PostRepository
    db   *gorm.DB
}

func NewPostService(repo *repository.PostRepository, db *gorm.DB) *PostService {
    return &PostService{
        repo: repo,
        db:   db,
    }
}

// CreatePost creates a new post
func (s *PostService) CreatePost(businessID uuid.UUID, content, imageURL string) (*models.Post, error) {
    post := &models.Post{
        BusinessID:  businessID,
        Content:     content,
        ImageURL:    imageURL,
        Likes:       0,
        Comments:    0,
        IsPublished: true,
    }

    if err := s.repo.CreatePost(post); err != nil {
        return nil, fmt.Errorf("failed to create post: %w", err)
    }

    return post, nil
}

// GetPostByID gets a post by ID
func (s *PostService) GetPostByID(id uuid.UUID) (*models.Post, error) {
    return s.repo.GetPostByID(id)
}

// GetPostsByBusinessID gets posts for a business
func (s *PostService) GetPostsByBusinessID(businessID uuid.UUID, page, pageSize int) ([]models.Post, int64, error) {
    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 20
    }

    offset := (page - 1) * pageSize
    return s.repo.GetPostsByBusinessID(businessID, pageSize, offset)
}

// GetFeedPosts gets posts from followed businesses
func (s *PostService) GetFeedPosts(businessIDs []uuid.UUID, page, pageSize int) ([]models.Post, int64, error) {
    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 20
    }

    offset := (page - 1) * pageSize
    return s.repo.GetFeedPosts(businessIDs, pageSize, offset)
}

// UpdatePost updates a post
func (s *PostService) UpdatePost(id uuid.UUID, content, imageURL string) (*models.Post, error) {
    post, err := s.repo.GetPostByID(id)
    if err != nil {
        return nil, err
    }
    if post == nil {
        return nil, errors.New("post not found")
    }

    if content != "" {
        post.Content = content
    }
    if imageURL != "" {
        post.ImageURL = imageURL
    }

    if err := s.repo.UpdatePost(post); err != nil {
        return nil, fmt.Errorf("failed to update post: %w", err)
    }

    return post, nil
}

// DeletePost deletes a post
func (s *PostService) DeletePost(id uuid.UUID) error {
    return s.repo.DeletePost(id)
}

// LikePost increments the like count
func (s *PostService) LikePost(id uuid.UUID) error {
    post, err := s.repo.GetPostByID(id)
    if err != nil {
        return err
    }
    if post == nil {
        return errors.New("post not found")
    }

    return s.repo.IncrementLikes(id)
}

// UnlikePost decrements the like count
func (s *PostService) UnlikePost(id uuid.UUID) error {
    post, err := s.repo.GetPostByID(id)
    if err != nil {
        return err
    }
    if post == nil {
        return errors.New("post not found")
    }

    if post.Likes > 0 {
        return s.repo.DecrementLikes(id)
    }
    return nil
}