package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/qper/hf/internal/domain"
)

var (
	ErrCategoryValidation = errors.New("invalid category payload")
	ErrCategoryConflict   = errors.New("category already exists")
	ErrCategoryNotFound   = errors.New("category not found")
	ErrCategoryForbidden  = errors.New("category access denied")
)

var hexColorPattern = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, userID string, req domain.CreateCategoryRequest) (*domain.Category, error)
	ListCategories(ctx context.Context, userID string) ([]domain.Category, error)
	GetCategoryByID(ctx context.Context, userID string, categoryID string) (*domain.Category, error)
	UpdateCategory(ctx context.Context, userID string, categoryID string, req domain.UpdateCategoryRequest) (*domain.Category, error)
	DeleteCategory(ctx context.Context, userID string, categoryID string) error
	ReorderCategories(ctx context.Context, userID string, ids []string) ([]domain.Category, error)
}

type CategoryService struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(ctx context.Context, userID string, req domain.CreateCategoryRequest) (*domain.Category, error) {
	if err := validateCreateCategory(req); err != nil {
		return nil, err
	}
	category, err := s.repo.CreateCategory(ctx, userID, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCategoryConflict
		}
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) List(ctx context.Context, userID string) ([]domain.Category, error) {
	return s.repo.ListCategories(ctx, userID)
}

func (s *CategoryService) GetByID(ctx context.Context, userID string, categoryID string) (*domain.Category, error) {
	category, err := s.repo.GetCategoryByID(ctx, userID, categoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}
	return category, nil
}

func (s *CategoryService) Update(ctx context.Context, userID string, categoryID string, req domain.UpdateCategoryRequest) (*domain.Category, error) {
	if err := validateUpdateCategory(req); err != nil {
		return nil, err
	}
	category, err := s.repo.UpdateCategory(ctx, userID, categoryID, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}
	return category, nil
}

func (s *CategoryService) Delete(ctx context.Context, userID string, categoryID string) error {
	err := s.repo.DeleteCategory(ctx, userID, categoryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrCategoryNotFound
		}
	}
	return err
}

func (s *CategoryService) Reorder(ctx context.Context, userID string, ids []string) ([]domain.Category, error) {
	if len(ids) == 0 {
		return nil, ErrCategoryValidation
	}
	categories, err := s.repo.ReorderCategories(ctx, userID, ids)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCategoryForbidden
		}
		return nil, err
	}
	if categories == nil {
		return nil, ErrCategoryForbidden
	}
	return categories, nil
}

func validateCreateCategory(req domain.CreateCategoryRequest) error {
	if strings.TrimSpace(req.Name) == "" || len(req.Name) > 100 {
		return ErrCategoryValidation
	}
	if !hexColorPattern.MatchString(req.Color) {
		return ErrCategoryValidation
	}
	if strings.TrimSpace(req.Icon) == "" {
		return ErrCategoryValidation
	}
	return nil
}

func validateUpdateCategory(req domain.UpdateCategoryRequest) error {
	if req.Name != nil && (strings.TrimSpace(*req.Name) == "" || len(*req.Name) > 100) {
		return ErrCategoryValidation
	}
	if req.Color != nil && !hexColorPattern.MatchString(*req.Color) {
		return ErrCategoryValidation
	}
	if req.Icon != nil && strings.TrimSpace(*req.Icon) == "" {
		return ErrCategoryValidation
	}
	return nil
}

func init() {
	_ = fmt.Sprintf("%v", domain.Category{})
}
