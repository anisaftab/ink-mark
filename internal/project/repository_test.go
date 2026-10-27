package project_test

import (
	"context"
	"testing"

	internaldb "ink-mark/internal/db"
	"ink-mark/internal/document"
	"ink-mark/internal/project"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const migrationsPath = "../../db/migrations"

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

func seedDoc(t *testing.T, db *gorm.DB, title, path string) document.Document {
	t.Helper()
	doc, err := document.NewRepository(db).Create(context.Background(), document.Input{Title: title, FilePath: path})
	require.NoError(t, err)
	return doc
}

func seedProject(t *testing.T, repo *project.Repository, name, folderPath string) project.Project {
	t.Helper()
	p, err := repo.Create(context.Background(), project.Input{Name: name, FolderPath: folderPath})
	require.NoError(t, err)
	return p
}

func TestRepository_Create(t *testing.T) {
	repo := project.NewRepository(setupTestDB(t))
	ctx := context.Background()

	tests := []struct {
		name    string
		input   project.Input
		want    func(t *testing.T, p project.Project)
		wantErr bool
	}{
		{
			name:  "success — name and folder path provided",
			input: project.Input{Name: "My Project", FolderPath: "/projects/my-project"},
			want: func(t *testing.T, p project.Project) {
				assert.Equal(t, "My Project", p.Name)
				assert.Equal(t, "/projects/my-project", p.FolderPath)
				assert.NotEmpty(t, p.ID)
				assert.False(t, p.CreatedAt.IsZero())
				assert.False(t, p.UpdatedAt.IsZero())
			},
		},
		{
			name:    "failure — empty name",
			input:   project.Input{Name: "", FolderPath: "/projects/x"},
			wantErr: true,
		},
		{
			name:    "failure — empty folder path",
			input:   project.Input{Name: "X", FolderPath: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := repo.Create(ctx, tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			tt.want(t, p)
		})
	}
}

func TestRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := project.NewRepository(db)
	ctx := context.Background()

	created := seedProject(t, repo, "Atlas", "/projects/atlas")

	tests := []struct {
		name    string
		id      string
		want    func(t *testing.T, p project.Project)
		wantErr error
	}{
		{
			name: "success — existing project",
			id:   created.ID,
			want: func(t *testing.T, p project.Project) {
				assert.Equal(t, created.ID, p.ID)
				assert.Equal(t, "Atlas", p.Name)
			},
		},
		{
			name:    "failure — non-existent id",
			id:      "00000000-0000-0000-0000-000000000000",
			wantErr: project.ErrNotFound,
		},
		{
			name:    "failure — empty id",
			id:      "",
			wantErr: project.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := repo.GetByID(ctx, tt.id)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			tt.want(t, p)
		})
	}
}

func TestRepository_List(t *testing.T) {
	t.Run("success — empty database returns empty slice", func(t *testing.T) {
		projects, err := project.NewRepository(setupTestDB(t)).List(context.Background())
		require.NoError(t, err)
		assert.Empty(t, projects)
	})

	t.Run("success — returns all projects", func(t *testing.T) {
		repo := project.NewRepository(setupTestDB(t))
		seedProject(t, repo, "First", "/first")
		seedProject(t, repo, "Second", "/second")

		projects, err := repo.List(context.Background())
		require.NoError(t, err)
		require.Len(t, projects, 2)
		names := []string{projects[0].Name, projects[1].Name}
		assert.Contains(t, names, "First")
		assert.Contains(t, names, "Second")
	})
}



func TestRepository_Delete(t *testing.T) {
	t.Run("success — project is deleted and no longer found", func(t *testing.T) {
		db := setupTestDB(t)
		repo := project.NewRepository(db)
		ctx := context.Background()

		p := seedProject(t, repo, "Delete Me", "/delete-me")
		require.NoError(t, repo.Delete(ctx, p.ID))
		_, err := repo.GetByID(ctx, p.ID)
		require.ErrorIs(t, err, project.ErrNotFound)
	})

	t.Run("success — deleting project sets document project_id to null", func(t *testing.T) {
		db := setupTestDB(t)
		repo := project.NewRepository(db)
		ctx := context.Background()

		p := seedProject(t, repo, "Cascade", "/cascade")
		doc := seedDoc(t, db, "Orphan", "/cascade/orphan.md")
		require.NoError(t, repo.AddDocument(ctx, p.ID, doc.ID))
		require.NoError(t, repo.Delete(ctx, p.ID))

		updated, err := document.NewRepository(db).GetByID(ctx, doc.ID)
		require.NoError(t, err)
		assert.Nil(t, updated.ProjectID)
	})

	t.Run("failure — non-existent id", func(t *testing.T) {
		err := project.NewRepository(setupTestDB(t)).Delete(context.Background(), "00000000-0000-0000-0000-000000000000")
		require.ErrorIs(t, err, project.ErrNotFound)
	})

	t.Run("failure — deleting an already deleted project", func(t *testing.T) {
		db := setupTestDB(t)
		repo := project.NewRepository(db)
		ctx := context.Background()

		p := seedProject(t, repo, "Double Delete", "/dd")
		require.NoError(t, repo.Delete(ctx, p.ID))
		require.ErrorIs(t, repo.Delete(ctx, p.ID), project.ErrNotFound)
	})
}

func TestRepository_DocumentMembership(t *testing.T) {
	t.Run("success — add document and list it", func(t *testing.T) {
		db := setupTestDB(t)
		repo := project.NewRepository(db)
		ctx := context.Background()

		p := seedProject(t, repo, "Handbook", "/handbook")
		doc := seedDoc(t, db, "Chapter 1", "/handbook/ch1.md")

		require.NoError(t, repo.AddDocument(ctx, p.ID, doc.ID))

		docs, err := repo.ListDocuments(ctx, p.ID)
		require.NoError(t, err)
		require.Len(t, docs, 1)
		assert.Equal(t, doc.ID, docs[0].ID)
		assert.Equal(t, &p.ID, docs[0].ProjectID)
	})

	t.Run("success — multiple documents all returned", func(t *testing.T) {
		db := setupTestDB(t)
		repo := project.NewRepository(db)
		ctx := context.Background()

		p := seedProject(t, repo, "Blog", "/blog")
		first := seedDoc(t, db, "Post A", "/blog/a.md")
		second := seedDoc(t, db, "Post B", "/blog/b.md")
		require.NoError(t, repo.AddDocument(ctx, p.ID, first.ID))
		require.NoError(t, repo.AddDocument(ctx, p.ID, second.ID))

		docs, err := repo.ListDocuments(ctx, p.ID)
		require.NoError(t, err)
		require.Len(t, docs, 2)
		ids := []string{docs[0].ID, docs[1].ID}
		assert.Contains(t, ids, first.ID)
		assert.Contains(t, ids, second.ID)
	})

	t.Run("success — remove document returns it to standalone", func(t *testing.T) {
		db := setupTestDB(t)
		repo := project.NewRepository(db)
		ctx := context.Background()

		p := seedProject(t, repo, "Temp", "/temp")
		doc := seedDoc(t, db, "Note", "/temp/note.md")
		require.NoError(t, repo.AddDocument(ctx, p.ID, doc.ID))
		require.NoError(t, repo.RemoveDocument(ctx, doc.ID))

		docs, err := repo.ListDocuments(ctx, p.ID)
		require.NoError(t, err)
		assert.Empty(t, docs)

		updated, err := document.NewRepository(db).GetByID(ctx, doc.ID)
		require.NoError(t, err)
		assert.Nil(t, updated.ProjectID)
	})

	t.Run("success — empty project returns empty slice", func(t *testing.T) {
		db := setupTestDB(t)
		repo := project.NewRepository(db)
		p := seedProject(t, repo, "Empty", "/empty")
		docs, err := repo.ListDocuments(context.Background(), p.ID)
		require.NoError(t, err)
		assert.Empty(t, docs)
	})

	t.Run("failure — add non-existent document returns document.ErrNotFound", func(t *testing.T) {
		db := setupTestDB(t)
		repo := project.NewRepository(db)
		p := seedProject(t, repo, "Project", "/project")
		err := repo.AddDocument(context.Background(), p.ID, "00000000-0000-0000-0000-000000000000")
		require.ErrorIs(t, err, document.ErrNotFound)
	})

	t.Run("failure — remove non-existent document returns document.ErrNotFound", func(t *testing.T) {
		err := project.NewRepository(setupTestDB(t)).RemoveDocument(context.Background(), "00000000-0000-0000-0000-000000000000")
		require.ErrorIs(t, err, document.ErrNotFound)
	})
}
