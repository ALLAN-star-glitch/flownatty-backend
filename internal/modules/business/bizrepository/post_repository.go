package bizrepository

import (
    "errors"
    "github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type PostRepository struct {
    db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
    return &PostRepository{db: db}
}

// CreatePost creates a new post
func (r *PostRepository) CreatePost(post *models.Post) error {
    return r.db.Create(post).Error
}

// GetPostByID gets a post by ID
func (r *PostRepository) GetPostByID(id uuid.UUID) (*models.Post, error) {
    var post models.Post
    err := r.db.Preload("Business").
        Where("id = ? AND is_published = ?", id, true).
        First(&post).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &post, nil
}

// GetPostsByBusinessID gets all posts for a business
func (r *PostRepository) GetPostsByBusinessID(businessID uuid.UUID, limit, offset int) ([]models.Post, int64, error) {
    var posts []models.Post
    var total int64
    
    db := r.db.Model(&models.Post{}).
        Where("business_id = ? AND is_published = ?", businessID, true)
    
    if err := db.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    err := db.Order("created_at DESC").
        Limit(limit).
        Offset(offset).
        Find(&posts).Error
    
    return posts, total, err
}

// GetFeedPosts gets posts from followed businesses
func (r *PostRepository) GetFeedPosts(businessIDs []uuid.UUID, limit, offset int) ([]models.Post, int64, error) {
    var posts []models.Post
    var total int64
    
    if len(businessIDs) == 0 {
        return posts, 0, nil
    }
    
    db := r.db.Model(&models.Post{}).
        Preload("Business").
        Where("business_id IN ? AND is_published = ?", businessIDs, true)
    
    if err := db.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    err := db.Order("created_at DESC").
        Limit(limit).
        Offset(offset).
        Find(&posts).Error
    
    return posts, total, err
}

// UpdatePost updates a post
func (r *PostRepository) UpdatePost(post *models.Post) error {
    return r.db.Save(post).Error
}

// DeletePost soft deletes a post
func (r *PostRepository) DeletePost(id uuid.UUID) error {
    return r.db.Delete(&models.Post{}, id).Error
}

// IncrementLikes increments the like count for a post
func (r *PostRepository) IncrementLikes(id uuid.UUID) error {
    return r.db.Model(&models.Post{}).
        Where("id = ?", id).
        UpdateColumn("likes", gorm.Expr("likes + ?", 1)).Error
}

// DecrementLikes decrements the like count for a post
func (r *PostRepository) DecrementLikes(id uuid.UUID) error {
    return r.db.Model(&models.Post{}).
        Where("id = ?", id).
        UpdateColumn("likes", gorm.Expr("likes - ?", 1)).Error
}