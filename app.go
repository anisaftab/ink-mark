package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"ink-mark/internal/db"
	"ink-mark/internal/document"
	"ink-mark/internal/logging"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/gorm"
)

// dataDir returns a reliable writable directory for app data.
// Uses %APPDATA%\ink-mark on Windows, ~/.ink-mark elsewhere.
func dataDir() string {
	if appData := os.Getenv("APPDATA"); appData != "" {
		dir := filepath.Join(appData, "ink-mark")
		_ = os.MkdirAll(dir, 0o755)
		return dir
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	dir := filepath.Join(home, ".ink-mark")
	_ = os.MkdirAll(dir, 0o755)
	return dir
}

// dbFilePath returns the absolute path to the SQLite database file.
func dbFilePath() string {
	return filepath.Join(dataDir(), "ink-mark.db")
}

// App is the main application struct bound to the Wails frontend.
type App struct {
	ctx        context.Context
	db         *gorm.DB
	documents  *document.Repository
	consoleLog *os.File // non-nil only when INKMARK_DEV_LOGS=1
}

// NewApp creates a new App application struct.
func NewApp() *App {
	return &App{}
}

// startup is called by Wails when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.consoleLog = logging.OpenLogFile("console")

	gormDB, err := db.Open(dbFilePath())
	if err != nil {
		runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Database Error",
			Message: "Failed to open database: " + err.Error(),
		})
		log.Fatalf("failed to open database: %v", err)
	}

	if err := runMigrations(gormDB); err != nil {
		runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Migration Error",
			Message: "Failed to run migrations: " + err.Error(),
		})
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
	if a.documents == nil {
		return DocumentDTO{}, fmt.Errorf("database not ready")
	}

	savePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "New Document",
		DefaultFilename: "Untitled.md",
		Filters:         []runtime.FileFilter{{DisplayName: "Markdown (*.md)", Pattern: "*.md"}},
	})
	if err != nil {
		return DocumentDTO{}, fmt.Errorf("save dialog: %w", err)
	}
	if savePath == "" {
		return DocumentDTO{}, nil
	}

	// Ensure the directory exists.
	if err := os.MkdirAll(filepath.Dir(savePath), 0o755); err != nil {
		return DocumentDTO{}, fmt.Errorf("create directory %q: %w", filepath.Dir(savePath), err)
	}

	// Create the empty .md file on disk.
	if err := os.WriteFile(savePath, []byte(""), 0o644); err != nil {
		return DocumentDTO{}, fmt.Errorf("create file %q: %w", savePath, err)
	}

	title := filepath.Base(savePath)
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
// Content is stored as-is (ProseKit JSON for now).
func (a *App) SaveDocument(id string, title string, content string) error {
	doc, err := a.documents.GetByID(a.ctx, id)
	if err != nil {
		return fmt.Errorf("get document %q: %w", id, err)
	}

	if err := os.WriteFile(doc.FilePath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write file %q: %w", doc.FilePath, err)
	}

	if title != "" && title != doc.Title {
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

// LogBrowserEvent is called by the frontend devlog patcher to forward
// console.error / console.warn / uncaught exceptions into logs/console.log.
// It is a no-op when INKMARK_DEV_LOGS is not "1".
func (a *App) LogBrowserEvent(level, message, stack string) {
	if a.consoleLog == nil {
		return
	}
	msg := message
	if stack != "" {
		msg = message + "\n" + stack
	}
	logging.WriteLine(a.consoleLog, level, "browser", msg)
}

