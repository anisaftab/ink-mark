package document

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Document represents a markdown document tracked in the database.
// The actual content lives on disk at FilePath.
type Document struct {
	ID        string    `gorm:"type:text;primaryKey"`
	Title     string    `gorm:"type:text;not null;default:'Untitled'"`
	FilePath  string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime"`
}

// TableName tells GORM to use the "documents" table.
func (Document) TableName() string {
	return "documents"
}

// BeforeCreate generates a UUID for new documents.
func (d *Document) BeforeCreate(_ *gorm.DB) error {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}

	return nil
}

// Input holds the fields for creating or updating a document.
type Input struct {
	Title    string
	FilePath string
}
