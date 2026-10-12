package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"ink-mark/internal/db"
	"ink-mark/internal/document"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/gorm"
)

const dbPath = "ink-mark.db"

// App is the main application struct bound to the Wails frontend.
type App struct {
	ctx       context.Context
	db        *gorm.DB
	documents *document.Repository
}

// NewApp creates a new App application struct.
func NewApp() *App {
	return &App{}
}

// startup is called by Wails when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	gormDB, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	if err := runMigrations(gormDB); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	a.db = gormDB
	a.documents = document.NewRepository(gormDB)
}

// DocumentDTO is the data transfer object returned to the frontend.
type DocumentDTO struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	FilePath string `json:"filePath"`
	Content  string `json:"content"` // file content read from disk
}

// NewDocument opens a save-file dialog to choose where the document lives,
// creates the empty .md file on disk, and records it in the database.
func (a *App) NewDocument() (DocumentDTO, error) {
	// Ask the user where to save the file first.
	savePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "New Document",
		DefaultFilename: "Untitled.md",
		Filters:         []runtime.FileFilter{{DisplayName: "Markdown (*.md)", Pattern: "*.md"}},
	})
	if err != nil {
		return DocumentDTO{}, fmt.Errorf("save dialog: %w", err)
	}
	if savePath == "" {
		// User cancelled.
		return DocumentDTO{}, nil
	}

	// Create the empty .md file on disk.
	if err := os.MkdirAll(filepath.Dir(savePath), 0o755); err != nil {
		return DocumentDTO{}, fmt.Errorf("create directory: %w", err)
	}
	if err := os.WriteFile(savePath, []byte(""), 0o644); err != nil {
		return DocumentDTO{}, fmt.Errorf("create file: %w", err)
	}

	title := filepath.Base(savePath)
	// Strip .md extension for the title.
	if ext := filepath.Ext(title); ext != "" {
		title = title[:len(title)-len(ext)]
	}

	doc, err := a.documents.Create(a.ctx, document.Input{
		Title:    title,
		FilePath: savePath,
	})
	if err != nil {
		return DocumentDTO{}, fmt.Errorf("create document record: %w", err)
	}

	return DocumentDTO{
		ID:       doc.ID,
		Title:    doc.Title,
		FilePath: doc.FilePath,
		Content:  "",
	}, nil
}

// SaveDocument writes content to the document's file on disk and updates the title in the DB.
func (a *App) SaveDocument(id string, title string, content string) error {
	doc, err := a.documents.GetByID(a.ctx, id)
	if err != nil {
		return fmt.Errorf("get document: %w", err)
	}

	if err := os.WriteFile(doc.FilePath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	if title != doc.Title {
		if _, err := a.documents.Update(a.ctx, id, document.Input{Title: title}); err != nil {
			return fmt.Errorf("update title: %w", err)
		}
	}

	return nil
}

// LoadDocument reads a document's content from disk.
func (a *App) LoadDocument(id string) (DocumentDTO, error) {
	doc, err := a.documents.GetByID(a.ctx, id)
	if err != nil {
		return DocumentDTO{}, fmt.Errorf("get document: %w", err)
	}

	data, err := os.ReadFile(doc.FilePath)
	if err != nil {
		return DocumentDTO{}, fmt.Errorf("read file %s: %w", doc.FilePath, err)
	}

	return DocumentDTO{
		ID:       doc.ID,
		Title:    doc.Title,
		FilePath: doc.FilePath,
		Content:  string(data),
	}, nil
}

