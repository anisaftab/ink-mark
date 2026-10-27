package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"ink-mark/internal/db"
	"ink-mark/internal/document"
	"ink-mark/internal/logging"
	"ink-mark/internal/project"

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
	projects   *project.Repository
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
	a.projects = project.NewRepository(gormDB)
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

// OpenDocument opens a file-picker dialog so the user can select an existing
// .md file. If the file is already tracked in the database its existing record
// is returned; otherwise a new record is created. Either way the file's content
// is read from disk and returned in the DTO.
func (a *App) OpenDocument() (DocumentDTO, error) {
	if a.documents == nil {
		return DocumentDTO{}, fmt.Errorf("database not ready")
	}

	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   "Open Document",
		Filters: []runtime.FileFilter{{DisplayName: "Markdown (*.md)", Pattern: "*.md"}},
	})
	if err != nil {
		return DocumentDTO{}, fmt.Errorf("open dialog: %w", err)
	}
	if filePath == "" {
		// User cancelled — return empty DTO with no error.
		return DocumentDTO{}, nil
	}

	// Re-use an existing record if this file is already known.
	doc, err := a.documents.GetByFilePath(a.ctx, filePath)
	if err != nil {
		if !errors.Is(err, document.ErrNotFound) {
			return DocumentDTO{}, fmt.Errorf("look up document by path: %w", err)
		}

		// File not yet in the database — register it now.
		title := filepath.Base(filePath)
		if ext := filepath.Ext(title); ext != "" {
			title = title[:len(title)-len(ext)]
		}

		doc, err = a.documents.Create(a.ctx, document.Input{
			Title:    title,
			FilePath: filePath,
		})
		if err != nil {
			return DocumentDTO{}, fmt.Errorf("create document record: %w", err)
		}
	}

	data, err := os.ReadFile(doc.FilePath)
	if err != nil {
		return DocumentDTO{}, fmt.Errorf("read file %q: %w", doc.FilePath, err)
	}

	return DocumentDTO{
		ID:       doc.ID,
		Title:    doc.Title,
		FilePath: doc.FilePath,
		Content:  string(data),
	}, nil
}

// ListDocuments returns all tracked documents ordered by most recently updated.
func (a *App) ListDocuments() ([]DocumentDTO, error) {
	if a.documents == nil {
		return nil, fmt.Errorf("database not ready")
	}

	docs, err := a.documents.List(a.ctx)
	if err != nil {
		return nil, fmt.Errorf("list documents: %w", err)
	}

	dtos := make([]DocumentDTO, len(docs))
	for i, doc := range docs {
		dtos[i] = DocumentDTO{
			ID:       doc.ID,
			Title:    doc.Title,
			FilePath: doc.FilePath,
			Content:  "", // content is loaded on demand via LoadDocument
		}
	}

	return dtos, nil
}

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

// ProjectDTO is the data transfer object for a project returned to the frontend.
type ProjectDTO struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	FolderPath string        `json:"folderPath"`
	Documents  []DocumentDTO `json:"documents"`
}

// NewProject opens a folder-picker dialog, creates a new project record in the
// database rooted at the chosen folder, and returns the empty project DTO.
func (a *App) NewProject() (ProjectDTO, error) {
	if a.projects == nil {
		return ProjectDTO{}, fmt.Errorf("database not ready")
	}

	folderPath, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Choose Project Folder",
	})
	if err != nil {
		return ProjectDTO{}, fmt.Errorf("folder dialog: %w", err)
	}
	if folderPath == "" {
		return ProjectDTO{}, nil // user cancelled
	}

	name := filepath.Base(folderPath)

	p, err := a.projects.Create(a.ctx, project.Input{
		Name:       name,
		FolderPath: folderPath,
	})
	if err != nil {
		return ProjectDTO{}, fmt.Errorf("create project: %w", err)
	}

	return ProjectDTO{
		ID:         p.ID,
		Name:       p.Name,
		FolderPath: p.FolderPath,
		Documents:  []DocumentDTO{},
	}, nil
}

// ListProjects returns all projects with their associated document metadata
// (content is not loaded — use LoadDocument for that).
func (a *App) ListProjects() ([]ProjectDTO, error) {
	if a.projects == nil {
		return nil, fmt.Errorf("database not ready")
	}

	projects, err := a.projects.List(a.ctx)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}

	dtos := make([]ProjectDTO, len(projects))
	for i, p := range projects {
		docs, err := a.projects.ListDocuments(a.ctx, p.ID)
		if err != nil {
			return nil, fmt.Errorf("list documents for project %s: %w", p.ID, err)
		}
		docDTOs := make([]DocumentDTO, len(docs))
		for j, doc := range docs {
			docDTOs[j] = DocumentDTO{
				ID:       doc.ID,
				Title:    doc.Title,
				FilePath: doc.FilePath,
				Content:  "",
			}
		}
		dtos[i] = ProjectDTO{
			ID:         p.ID,
			Name:       p.Name,
			FolderPath: p.FolderPath,
			Documents:  docDTOs,
		}
	}

	return dtos, nil
}

// GetProjectDocuments returns all document metadata for a single project.
func (a *App) GetProjectDocuments(projectID string) ([]DocumentDTO, error) {
	if a.projects == nil {
		return nil, fmt.Errorf("database not ready")
	}

	docs, err := a.projects.ListDocuments(a.ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("get project documents: %w", err)
	}

	dtos := make([]DocumentDTO, len(docs))
	for i, doc := range docs {
		dtos[i] = DocumentDTO{
			ID:       doc.ID,
			Title:    doc.Title,
			FilePath: doc.FilePath,
			Content:  "",
		}
	}

	return dtos, nil
}

// AddDocumentToProject assigns an existing document to a project.
func (a *App) AddDocumentToProject(projectID, documentID string) error {
	if a.projects == nil {
		return fmt.Errorf("database not ready")
	}

	// Verify the project exists before attempting the assignment.
	if _, err := a.projects.GetByID(a.ctx, projectID); err != nil {
		if errors.Is(err, project.ErrNotFound) {
			return fmt.Errorf("project not found")
		}
		return fmt.Errorf("get project: %w", err)
	}

	if err := a.projects.AddDocument(a.ctx, projectID, documentID); err != nil {
		return fmt.Errorf("add document to project: %w", err)
	}

	return nil
}

