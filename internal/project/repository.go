package project

import (
	"context"
	"errors"
	"fmt"

	"ink-mark/internal/document"

	"gorm.io/gorm"
)

// ErrNotFound is returned when a project does not exist.
var ErrNotFound = errors.New("project not found")

// Repository handles all database operations for projects.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new Repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create inserts a new project record.
func (r *Repository) Create(ctx context.Context, input Input) (Project, error) {
	if input.Name == "" {
		return Project{}, fmt.Errorf("create project: name is required")
	}
	if input.FolderPath == "" {
		return Project{}, fmt.Errorf("create project: folder path is required")
	}

	p := Project{
		Name:       input.Name,
		FolderPath: input.FolderPath,
	}

	if err := r.db.WithContext(ctx).Create(&p).Error; err != nil {
		return Project{}, fmt.Errorf("create project: %w", err)
	}

	return p, nil
}

// GetByID returns a single project by its ID.
// Returns ErrNotFound if no project exists with that ID.
func (r *Repository) GetByID(ctx context.Context, id string) (Project, error) {
	var p Project

	result := r.db.WithContext(ctx).First(&p, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Project{}, ErrNotFound
		}
		return Project{}, fmt.Errorf("get project %s: %w", id, result.Error)
	}

	return p, nil
}

// GetByName returns the first project whose name matches exactly (case-insensitive).
// Returns ErrNotFound if no project exists with that name.
func (r *Repository) GetByName(ctx context.Context, name string) (Project, error) {
	var p Project

	result := r.db.WithContext(ctx).
		Where("lower(name) = lower(?)", name).
		First(&p)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Project{}, ErrNotFound
		}
		return Project{}, fmt.Errorf("get project by name %q: %w", name, result.Error)
	}

	return p, nil
}

// List returns all projects ordered by most recently updated first.
func (r *Repository) List(ctx context.Context) ([]Project, error) {
	var projects []Project

	if err := r.db.WithContext(ctx).Order("updated_at desc").Find(&projects).Error; err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}

	return projects, nil
}

// Delete removes a project record by ID.
// Documents that belonged to the project have their project_id set to NULL
// before the project row is deleted, ensuring the disassociation works
// regardless of whether the SQLite driver has FK enforcement enabled.
// Returns ErrNotFound if no project exists with that ID.
func (r *Repository) Delete(ctx context.Context, id string) error {
	// First verify the project exists.
	if _, err := r.GetByID(ctx, id); err != nil {
		return err
	}

	// Disassociate all documents from this project.
	if err := r.db.WithContext(ctx).
		Model(&document.Document{}).
		Where("project_id = ?", id).
		Update("project_id", nil).Error; err != nil {
		return fmt.Errorf("disassociate documents from project %s: %w", id, err)
	}

	// Delete the project row.
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&Project{}).Error; err != nil {
		return fmt.Errorf("delete project %s: %w", id, err)
	}

	return nil
}

// ListDocuments returns all documents that belong to the given project,
// ordered by most recently updated first.
func (r *Repository) ListDocuments(ctx context.Context, projectID string) ([]document.Document, error) {
	var docs []document.Document

	result := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Order("updated_at desc").
		Find(&docs)

	if result.Error != nil {
		return nil, fmt.Errorf("list documents for project %s: %w", projectID, result.Error)
	}

	return docs, nil
}

// AddDocument assigns an existing document to this project.
// Returns document.ErrNotFound if the document does not exist.
func (r *Repository) AddDocument(ctx context.Context, projectID, documentID string) error {
	result := r.db.WithContext(ctx).
		Model(&document.Document{}).
		Where("id = ?", documentID).
		Update("project_id", projectID)

	if result.Error != nil {
		return fmt.Errorf("add document %s to project %s: %w", documentID, projectID, result.Error)
	}
	if result.RowsAffected == 0 {
		return document.ErrNotFound
	}

	return nil
}

// RemoveDocument removes the project association from a document (sets project_id to NULL).
// Returns document.ErrNotFound if the document does not exist.
func (r *Repository) RemoveDocument(ctx context.Context, documentID string) error {
	result := r.db.WithContext(ctx).
		Model(&document.Document{}).
		Where("id = ?", documentID).
		Update("project_id", nil)

	if result.Error != nil {
		return fmt.Errorf("remove document %s from project: %w", documentID, result.Error)
	}
	if result.RowsAffected == 0 {
		return document.ErrNotFound
	}

	return nil
}
