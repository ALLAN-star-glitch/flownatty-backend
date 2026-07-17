package bizservice

import (
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/business/bizrepository"
	"github.com/google/uuid"
)

type CategoryService struct {
    repo *bizrepository.CategoryRepository
}

func NewCategoryService(repo *bizrepository.CategoryRepository) *CategoryService {
    return &CategoryService{repo: repo}
}

// GetAllCategories gets all active categories
func (s *CategoryService) GetAllCategories() ([]models.ProductServiceCategory, error) {
    return s.repo.GetAllCategories()
}

// GetCategoryByID gets a category by ID
func (s *CategoryService) GetCategoryByID(id uuid.UUID) (*models.ProductServiceCategory, error) {
    return s.repo.GetCategoryByID(id)
}

// GetCategoryBySlug gets a category by slug
func (s *CategoryService) GetCategoryBySlug(slug string) (*models.ProductServiceCategory, error) {
    return s.repo.GetCategoryBySlug(slug)
}