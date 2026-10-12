package document

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// isInvalidUUIDError returns true if the error is a Postgres invalid UUID syntax error.
func isInvalidUUIDError(err error) bool {
	return strings.Contains(err.Error(), "invalid input syntax for type uuid")
}

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

// Create inserts a new document and returns it with all database-generated fields populated.
func (r *Repository) Create(ctx context.Context, input Input) (Document, error) {
	doc := Document{
		Title:   input.Title,
		Content: input.Content,
	}

	if doc.Title == "" {
		doc.Title = "Untitled"
	}

	if len(doc.Content) == 0 {
		doc.Content = json.RawMessage(`{}`)
	}

	result := r.db.WithContext(ctx).Create(&doc)
	if result.Error != nil {
		return Document{}, fmt.Errorf("create document: %w", result.Error)
	}

	return doc, nil
}

// GetByID returns a single document by its UUID.
// Returns ErrNotFound if no document exists with that ID or the ID is malformed.
func (r *Repository) GetByID(ctx context.Context, id string) (Document, error) {
	var doc Document

	result := r.db.WithContext(ctx).First(&doc, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Document{}, ErrNotFound
		}
		// Postgres returns a syntax error for malformed UUIDs — treat as not found.
		if isInvalidUUIDError(result.Error) {
			return Document{}, ErrNotFound
		}
		return Document{}, fmt.Errorf("get document %s: %w", id, result.Error)
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

// Update modifies the title and/or content of an existing document.
// Returns ErrNotFound if no document exists with that ID.
func (r *Repository) Update(ctx context.Context, id string, input Input) (Document, error) {
	updates := map[string]any{}

	if input.Title != "" {
		updates["title"] = input.Title
	}
	if len(input.Content) > 0 {
		updates["content"] = input.Content
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

// Delete removes a document by its UUID.
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
