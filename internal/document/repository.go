package document

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// ErrNotFound is returned when a document does not exist.
var ErrNotFound = errors.New("document not found")

// Repository handles all database operations for documents.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new Repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create inserts a new document record.
func (r *Repository) Create(ctx context.Context, input Input) (Document, error) {
	if input.FilePath == "" {
		return Document{}, fmt.Errorf("create document: file path is required")
	}

	doc := Document{
		Title:    input.Title,
		FilePath: input.FilePath,
	}

	if doc.Title == "" {
		doc.Title = "Untitled"
	}

	result := r.db.WithContext(ctx).Create(&doc)
	if result.Error != nil {
		return Document{}, fmt.Errorf("create document: %w", result.Error)
	}

	return doc, nil
}

// GetByID returns a single document by its ID.
// Returns ErrNotFound if no document exists with that ID.
func (r *Repository) GetByID(ctx context.Context, id string) (Document, error) {
	var doc Document

	result := r.db.WithContext(ctx).First(&doc, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Document{}, ErrNotFound
		}
		return Document{}, fmt.Errorf("get document %s: %w", id, result.Error)
	}

	return doc, nil
}

// GetByFilePath returns the document whose file_path matches exactly.
// Returns ErrNotFound if no document exists with that path.
func (r *Repository) GetByFilePath(ctx context.Context, filePath string) (Document, error) {
	var doc Document

	result := r.db.WithContext(ctx).Where("file_path = ?", filePath).First(&doc)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Document{}, ErrNotFound
		}
		return Document{}, fmt.Errorf("get document by path %q: %w", filePath, result.Error)
	}

	return doc, nil
}

// GetByTitle returns the first document matching the given title (case-insensitive).
// Returns ErrNotFound if no document exists with that title.
func (r *Repository) GetByTitle(ctx context.Context, title string) (Document, error) {
	var doc Document

	result := r.db.WithContext(ctx).
		Where("lower(title) = lower(?)", title).
		First(&doc)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Document{}, ErrNotFound
		}
		return Document{}, fmt.Errorf("get document by title %q: %w", title, result.Error)
	}

	return doc, nil
}

// List returns all documents ordered by most recently updated first.
func (r *Repository) List(ctx context.Context) ([]Document, error) {
	var docs []Document

	result := r.db.WithContext(ctx).Order("updated_at desc").Find(&docs)
	if result.Error != nil {
		return nil, fmt.Errorf("list documents: %w", result.Error)
	}

	return docs, nil
}

// Update modifies the title and/or file path of an existing document.
// Returns ErrNotFound if no document exists with that ID.
func (r *Repository) Update(ctx context.Context, id string, input Input) (Document, error) {
	updates := map[string]any{}

	if input.Title != "" {
		updates["title"] = input.Title
	}
	if input.FilePath != "" {
		updates["file_path"] = input.FilePath
	}

	result := r.db.WithContext(ctx).
		Model(&Document{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return Document{}, fmt.Errorf("update document %s: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return Document{}, ErrNotFound
	}

	return r.GetByID(ctx, id)
}

// Delete removes a document record by ID.
// Returns ErrNotFound if no document exists with that ID.
func (r *Repository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&Document{})
	if result.Error != nil {
		return fmt.Errorf("delete document %s: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

