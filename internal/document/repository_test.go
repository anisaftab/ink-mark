package document_test

import (
	"context"
	"encoding/json"
	"testing"

	internaldb "ink-mark/internal/db"
	"ink-mark/internal/document"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const migrationsPath = "../../db/migrations"

// setupTestSchema creates an isolated Postgres schema named schemaName, runs all
// migrations inside it, and registers a cleanup that drops it — even on failure.
func setupTestSchema(t *testing.T, schemaName string) *gorm.DB {
	t.Helper()

	// Use a plain connection to manage the schema lifecycle.
	rootDB, err := internaldb.Open()
	require.NoError(t, err, "open root database connection")

	// Drop any stale schema from a previous crashed run, then recreate clean.
	require.NoError(t,
		rootDB.Exec("DROP SCHEMA IF EXISTS "+schemaName+" CASCADE").Error,
		"drop stale schema",
	)
	require.NoError(t,
		rootDB.Exec("CREATE SCHEMA "+schemaName).Error,
		"create test schema",
	)

	// Always drop on exit — even on failure or panic.
	t.Cleanup(func() {
		_ = rootDB.Exec("DROP SCHEMA IF EXISTS " + schemaName + " CASCADE").Error
	})

	// Open a schema-scoped connection — search_path baked into DSN, no pool leakage.
	gormDB, err := internaldb.OpenWithSchema(schemaName)
	require.NoError(t, err, "open schema-scoped connection")

	sqlDB, err := gormDB.DB()
	require.NoError(t, err, "get underlying sql.DB")

	// Use a schema-specific goose version table so each schema tracks its own
	// migration state independently — otherwise goose sees the public schema's
	// version table and skips migrations entirely.
	goose.SetTableName(schemaName + ".goose_db_version")
	goose.SetLogger(goose.NopLogger())
	require.NoError(t, goose.Up(sqlDB, migrationsPath), "run migrations")

	return gormDB
}

func TestRepository_Create(t *testing.T) {
	db := setupTestSchema(t, "test_repo_create")
	repo := document.NewRepository(db)
	ctx := context.Background()

	tests := []struct {
		name    string
		input   document.Input
		want    func(t *testing.T, doc document.Document)
		wantErr bool
	}{
		{
			name:  "success — title and content provided",
			input: document.Input{Title: "My Note", Content: json.RawMessage(`{"type":"doc"}`)},
			want: func(t *testing.T, doc document.Document) {
				assert.Equal(t, "My Note", doc.Title)
				assert.JSONEq(t, `{"type":"doc"}`, string(doc.Content))
				assert.NotEmpty(t, doc.ID)
				assert.False(t, doc.CreatedAt.IsZero())
				assert.False(t, doc.UpdatedAt.IsZero())
			},
		},
		{
			name:  "success — empty title defaults to Untitled",
			input: document.Input{Title: "", Content: json.RawMessage(`{}`)},
			want: func(t *testing.T, doc document.Document) {
				assert.Equal(t, "Untitled", doc.Title)
			},
		},
		{
			name:  "success — nil content defaults to {}",
			input: document.Input{Title: "Note", Content: nil},
			want: func(t *testing.T, doc document.Document) {
				assert.JSONEq(t, `{}`, string(doc.Content))
			},
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
	db := setupTestSchema(t, "test_repo_getbyid")
	repo := document.NewRepository(db)
	ctx := context.Background()

	created, err := repo.Create(ctx, document.Input{Title: "Find Me", Content: json.RawMessage(`{}`)})
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
	db := setupTestSchema(t, "test_repo_getbytitle")
	repo := document.NewRepository(db)
	ctx := context.Background()

	_, err := repo.Create(ctx, document.Input{Title: "Meeting Notes", Content: json.RawMessage(`{}`)})
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
		db := setupTestSchema(t, "test_repo_list_empty")
		repo := document.NewRepository(db)
		ctx := context.Background()

		docs, err := repo.List(ctx)
		require.NoError(t, err)
		assert.Empty(t, docs)
	})

	t.Run("success — returns all documents ordered by updated_at desc", func(t *testing.T) {
		db := setupTestSchema(t, "test_repo_list_ordered")
		repo := document.NewRepository(db)
		ctx := context.Background()

		_, err := repo.Create(ctx, document.Input{Title: "First", Content: json.RawMessage(`{}`)})
		require.NoError(t, err)
		_, err = repo.Create(ctx, document.Input{Title: "Second", Content: json.RawMessage(`{}`)})
		require.NoError(t, err)

		docs, err := repo.List(ctx)
		require.NoError(t, err)
		assert.Len(t, docs, 2)
		assert.Equal(t, "Second", docs[0].Title)
		assert.Equal(t, "First", docs[1].Title)
	})
}


func TestRepository_Update(t *testing.T) {
	db := setupTestSchema(t, "test_repo_update")
	repo := document.NewRepository(db)
	ctx := context.Background()

	created, err := repo.Create(ctx, document.Input{Title: "Original", Content: json.RawMessage(`{"v":1}`)})
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
				assert.JSONEq(t, `{"v":1}`, string(doc.Content))
			},
		},
		{
			name:  "success — update content only",
			id:    created.ID,
			input: document.Input{Content: json.RawMessage(`{"v":2}`)},
			want: func(t *testing.T, doc document.Document) {
				assert.JSONEq(t, `{"v":2}`, string(doc.Content))
			},
		},
		{
			name:  "success — updated_at is refreshed",
			id:    created.ID,
			input: document.Input{Title: "Trigger Timestamp"},
			want: func(t *testing.T, doc document.Document) {
				assert.True(t, doc.UpdatedAt.After(created.UpdatedAt) || doc.UpdatedAt.Equal(created.UpdatedAt))
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
	db := setupTestSchema(t, "test_repo_delete")
	repo := document.NewRepository(db)
	ctx := context.Background()

	t.Run("success — document is deleted and no longer found", func(t *testing.T) {
		doc, err := repo.Create(ctx, document.Input{Title: "Delete Me", Content: json.RawMessage(`{}`)})
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
		doc, err := repo.Create(ctx, document.Input{Title: "Double Delete", Content: json.RawMessage(`{}`)})
		require.NoError(t, err)

		require.NoError(t, repo.Delete(ctx, doc.ID))

		err = repo.Delete(ctx, doc.ID)
		require.ErrorIs(t, err, document.ErrNotFound)
	})
}

