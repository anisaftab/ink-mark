package document

import (
	"encoding/json"
	"time"
)

// Document represents a single markdown document stored in the database.
type Document struct {
	ID        string          `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title     string          `gorm:"type:text;not null;default:'Untitled'"`
	Content   json.RawMessage `gorm:"type:jsonb;not null;default:'{}'"`
	CreatedAt time.Time       `gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time       `gorm:"not null;autoUpdateTime"`
}

// TableName tells GORM to use the "documents" table.
func (Document) TableName() string {
	return "documents"
}

// Input holds the fields for creating or updating a document.
type Input struct {
	Title   string
	Content json.RawMessage
}
