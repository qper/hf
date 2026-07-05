package service

import (
	"context"
	"errors"
	"testing"

	"github.com/qper/hf/internal/domain"
)

type stubCategoryRepository struct {
	created *domain.Category
	got     *domain.Category
}

func (s stubCategoryRepository) CreateCategory(ctx context.Context, userID string, req domain.CreateCategoryRequest) (*domain.Category, error) {
	return s.created, nil
}

func (s stubCategoryRepository) ListCategories(ctx context.Context, userID string) ([]domain.Category, error) {
	return nil, nil
}

func (s stubCategoryRepository) GetCategoryByID(ctx context.Context, userID string, categoryID string) (*domain.Category, error) {
	return s.got, nil
}

func (s stubCategoryRepository) UpdateCategory(ctx context.Context, userID string, categoryID string, req domain.UpdateCategoryRequest) (*domain.Category, error) {
	return s.got, nil
}

func (s stubCategoryRepository) DeleteCategory(ctx context.Context, userID string, categoryID string) error {
	return nil
}

func (s stubCategoryRepository) ReorderCategories(ctx context.Context, userID string, ids []string) ([]domain.Category, error) {
	return nil, nil
}

func TestCreateCategoryRejectsInvalidPayload(t *testing.T) {
	svc := NewCategoryService(stubCategoryRepository{})

	_, err := svc.Create(context.Background(), "u1", domain.CreateCategoryRequest{Name: "", Color: "#123456", Icon: "home"})
	if !errors.Is(err, ErrCategoryValidation) {
		t.Fatalf("expected ErrCategoryValidation, got %v", err)
	}
}

func TestCreateCategoryRejectsInvalidColor(t *testing.T) {
	svc := NewCategoryService(stubCategoryRepository{})

	_, err := svc.Create(context.Background(), "u1", domain.CreateCategoryRequest{Name: "Work", Color: "blue", Icon: "home"})
	if !errors.Is(err, ErrCategoryValidation) {
		t.Fatalf("expected ErrCategoryValidation, got %v", err)
	}
}

func TestGetCategoryReturnsErrNotFoundForMissingCategory(t *testing.T) {
	svc := NewCategoryService(stubCategoryRepository{})

	_, err := svc.GetByID(context.Background(), "u1", "missing")
	if !errors.Is(err, ErrCategoryNotFound) {
		t.Fatalf("expected ErrCategoryNotFound, got %v", err)
	}
}
