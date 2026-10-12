package db

import (
"fmt"

"github.com/glebarez/sqlite"
"gorm.io/gorm"
"gorm.io/gorm/logger"
)

// Open opens a GORM SQLite connection at the given file path.
// The file is created if it does not exist.
func Open(path string) (*gorm.DB, error) {
db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
Logger: logger.Default.LogMode(logger.Silent),
})
if err != nil {
return nil, fmt.Errorf("open database %s: %w", path, err)
}
return db, nil
}

// OpenMemory opens an in-memory SQLite database.
// Used in tests — each call returns a completely isolated instance.
func OpenMemory() (*gorm.DB, error) {
db, err := gorm.Open(sqlite.Open("file::memory:?cache=private"), &gorm.Config{
Logger: logger.Default.LogMode(logger.Silent),
})
if err != nil {
return nil, fmt.Errorf("open in-memory database: %w", err)
}
return db, nil
}
