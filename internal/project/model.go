package project

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Project represents a named collection of documents that maps to a folder on disk.
type Project struct {
	ID         string    `gorm:"type:text;primaryKey"`
	Name       string    `gorm:"type:text;not null"`
	FolderPath string    `gorm:"type:text;not null"`
	CreatedAt  time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"not null;autoUpdateTime"`
}

// TableName tells GORM to use the "projects" table.
func (Project) TableName() string {
	return "projects"
}

// BeforeCreate generates a UUID for new projects.
func (p *Project) BeforeCreate(_ *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}

	return nil
}

// Input holds the fields for creating or updating a project.
type Input struct {
	Name       string
	FolderPath string
}
