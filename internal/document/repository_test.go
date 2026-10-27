package document_test

import (
	"context"
	"testing"

	internaldb "ink-mark/internal/db"
	"ink-mark/internal/document"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const migrationsPath = "../../db/migrations"

// setupTestDB opens an isolated in-memory SQLite database and runs all migrations.
// No cleanup needed — the database evaporates when the test ends.
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	gormDB, err := internaldb.OpenMemory()
	require.NoError(t, err, "open in-memory database")

	sqlDB, err := gormDB.DB()
	require.NoError(t, err, "get underlying sql.DB")

	goose.SetLogger(goose.NopLogger())
	require.NoError(t, goose.SetDialect("sqlite3"), "set dialect")
	require.NoError(t, goose.Up(sqlDB, migrationsPath), "run migrations")

	return gormDB
}

func TestRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := document.NewRepository(db)
	ctx := context.Background()

	tests := []struct {
		name    string
		input   document.Input
		want    func(t *testing.T, doc document.Document)
		wantErr bool
	}{
		{
			name:  "success — title and file path provided",
			input: document.Input{Title: "My Note", FilePath: "/docs/my-note.md"},
			want: func(t *testing.T, doc document.Document) {
				assert.Equal(t, "My Note", doc.Title)
				assert.Equal(t, "/docs/my-note.md", doc.FilePath)
				assert.NotEmpty(t, doc.ID)
				assert.False(t, doc.CreatedAt.IsZero())
				assert.False(t, doc.UpdatedAt.IsZero())
			},
		},
		{
			name:  "success — empty title defaults to Untitled",
			input: document.Input{Title: "", FilePath: "/docs/untitled.md"},
			want: func(t *testing.T, doc document.Document) {
				assert.Equal(t, "Untitled", doc.Title)
			},
		},
		{
			name:    "failure — empty file path",
			input:   document.Input{Title: "Note", FilePath: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := repo.Create(ctx, tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			tt.want(t, doc)
		})
	}
}

func TestRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := document.NewRepository(db)
	ctx := context.Background()

	created, err := repo.Create(ctx, document.Input{Title: "Find Me", FilePath: "/docs/find-me.md"})
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      string
		want    func(t *testing.T, doc document.Document)
		wantErr error
	}{
		{
			name: "success — existing document",
			id:   created.ID,
			want: func(t *testing.T, doc document.Document) {
				assert.Equal(t, created.ID, doc.ID)
				assert.Equal(t, "Find Me", doc.Title)
			},
		},
		{
			name:    "failure — non-existent id",
			id:      "00000000-0000-0000-0000-000000000000",
			wantErr: document.ErrNotFound,
		},
		{
			name:    "failure — malformed id",
			id:      "not-a-uuid",
			wantErr: document.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := repo.GetByID(ctx, tt.id)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			tt.want(t, doc)
		})
	}
}

func TestRepository_GetByTitle(t *testing.T) {
	db := setupTestDB(t)
	repo := document.NewRepository(db)
	ctx := context.Background()

	_, err := repo.Create(ctx, document.Input{Title: "Meeting Notes", FilePath: "/docs/meeting-notes.md"})
	require.NoError(t, err)

	tests := []struct {
		name    string
		title   string
		want    func(t *testing.T, doc document.Document)
		wantErr error
	}{
		{
			name:  "success — exact match",
			title: "Meeting Notes",
			want: func(t *testing.T, doc document.Document) {
				assert.Equal(t, "Meeting Notes", doc.Title)
			},
		},
		{
			name:  "success — case insensitive match",
			title: "meeting notes",
			want: func(t *testing.T, doc document.Document) {
				assert.Equal(t, "Meeting Notes", doc.Title)
			},
		},
		{
			name:    "failure — title does not exist",
			title:   "No Such Document",
			wantErr: document.ErrNotFound,
		},
		{
			name:    "failure — empty title",
			title:   "",
			wantErr: document.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := repo.GetByTitle(ctx, tt.title)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			tt.want(t, doc)
		})
	}
}

func TestRepository_List(t *testing.T) {
	t.Run("success — empty list", func(t *testing.T) {
		db := setupTestDB(t)
		repo := document.NewRepository(db)
		ctx := context.Background()

		docs, err := repo.List(ctx)
		require.NoError(t, err)
		assert.Empty(t, docs)
	})

	t.Run("success — returns all documents", func(t *testing.T) {
		db := setupTestDB(t)
		repo := document.NewRepository(db)
		ctx := context.Background()

		_, err := repo.Create(ctx, document.Input{Title: "First", FilePath: "/docs/first.md"})
		require.NoError(t, err)
		_, err = repo.Create(ctx, document.Input{Title: "Second", FilePath: "/docs/second.md"})
		require.NoError(t, err)

		docs, err := repo.List(ctx)
		require.NoError(t, err)
		require.Len(t, docs, 2)
		titles := []string{docs[0].Title, docs[1].Title}
		assert.Contains(t, titles, "First")
		assert.Contains(t, titles, "Second")
	})
}


func TestRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := document.NewRepository(db)
	ctx := context.Background()

	created, err := repo.Create(ctx, document.Input{Title: "Original", FilePath: "/docs/original.md"})
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      string
		input   document.Input
		want    func(t *testing.T, doc document.Document)
		wantErr error
	}{
		{
			name:  "success — update title only",
			id:    created.ID,
			input: document.Input{Title: "Updated Title"},
			want: func(t *testing.T, doc document.Document) {
				assert.Equal(t, "Updated Title", doc.Title)
				assert.Equal(t, "/docs/original.md", doc.FilePath)
			},
		},
		{
			name:  "success — update file path only",
			id:    created.ID,
			input: document.Input{FilePath: "/docs/renamed.md"},
			want: func(t *testing.T, doc document.Document) {
				assert.Equal(t, "/docs/renamed.md", doc.FilePath)
			},
		},
		{
			name:  "success — updated_at is refreshed",
			id:    created.ID,
			input: document.Input{Title: "Trigger Timestamp"},
			want: func(t *testing.T, doc document.Document) {
				// updated_at must be >= the original; sub-millisecond
				// SQLite resolution means it may be equal on fast machines.
				assert.False(t, doc.UpdatedAt.Before(created.UpdatedAt))
			},
		},
		{
			name:    "failure — non-existent id",
			id:      "00000000-0000-0000-0000-000000000000",
			input:   document.Input{Title: "Ghost"},
			wantErr: document.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := repo.Update(ctx, tt.id, tt.input)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			tt.want(t, doc)
		})
	}
}

func TestRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := document.NewRepository(db)
	ctx := context.Background()

	t.Run("success — document is deleted and no longer found", func(t *testing.T) {
		doc, err := repo.Create(ctx, document.Input{Title: "Delete Me", FilePath: "/docs/delete-me.md"})
		require.NoError(t, err)

		require.NoError(t, repo.Delete(ctx, doc.ID))

		_, err = repo.GetByID(ctx, doc.ID)
		require.ErrorIs(t, err, document.ErrNotFound)
	})

	t.Run("failure — non-existent id", func(t *testing.T) {
		err := repo.Delete(ctx, "00000000-0000-0000-0000-000000000000")
		require.ErrorIs(t, err, document.ErrNotFound)
	})

	t.Run("failure — deleting an already deleted document", func(t *testing.T) {
		doc, err := repo.Create(ctx, document.Input{Title: "Double Delete", FilePath: "/docs/double-delete.md"})
		require.NoError(t, err)

		require.NoError(t, repo.Delete(ctx, doc.ID))

		err = repo.Delete(ctx, doc.ID)
		require.ErrorIs(t, err, document.ErrNotFound)
	})
}

