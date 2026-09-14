package service

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// BackupService handles SQLite file export and import.
type BackupService struct {
	ctx      context.Context
	db       *sql.DB
	dbPath   string
	reopenDB func() error
}

// Init sets the runtime dependencies after Wails startup.
func (s *BackupService) Init(ctx context.Context, db *sql.DB, dbPath string, reopenDB func() error) {
	s.ctx = ctx
	s.db = db
	s.dbPath = dbPath
	s.reopenDB = reopenDB
}

// Export opens a save dialog and copies the SQLite database file to the selected path.
func (s *BackupService) Export() error {
	defaultDir := filepath.Dir(s.dbPath)
	filename := fmt.Sprintf("soditapp_backup_%s.db", time.Now().Format("20060102_150405"))

	path, err := runtime.SaveFileDialog(s.ctx, runtime.SaveDialogOptions{
		DefaultFilename: filename,
		DefaultDirectory: defaultDir,
		Filters: []runtime.FileFilter{
			{DisplayName: "SQLite Database", Pattern: "*.db"},
		},
	})
	if err != nil {
		return fmt.Errorf("opening save dialog: %w", err)
	}
	if path == "" {
		return nil // cancelled
	}

	src, err := os.Open(s.dbPath)
	if err != nil {
		return fmt.Errorf("opening source database: %w", err)
	}
	defer src.Close()

	// Copy WAL/SHM files if they exist (WAL mode)
	walPath := s.dbPath + "-wal"
	shmPath := s.dbPath + "-shm"

	if err := copyFile(s.dbPath, path); err != nil {
		return fmt.Errorf("copying database: %w", err)
	}

	// Best-effort copy of WAL and SHM files
	copyFile(walPath, path+"-wal")
	copyFile(shmPath, path+"-shm")

	return nil
}

// Import opens a file dialog, replaces the current database with the selected one,
// and reopens the connection.
func (s *BackupService) Import() error {
	path, err := runtime.OpenFileDialog(s.ctx, runtime.OpenDialogOptions{
		Filters: []runtime.FileFilter{
			{DisplayName: "SQLite Database", Pattern: "*.db"},
		},
	})
	if err != nil {
		return fmt.Errorf("opening file dialog: %w", err)
	}
	if path == "" {
		return nil // cancelled
	}

	// Close current connection
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("closing current database: %w", err)
	}

	// Replace current DB file with the backup
	if err := copyFile(path, s.dbPath); err != nil {
		// Try to reopen the old database even if copy failed
		s.reopenDB()
		return fmt.Errorf("copying backup file: %w", err)
	}

	// Copy WAL/SHM from backup if they exist
	copyFile(path+"-wal", s.dbPath+"-wal")
	copyFile(path+"-shm", s.dbPath+"-shm")

	// Reopen connection
	if err := s.reopenDB(); err != nil {
		return fmt.Errorf("reopening database: %w", err)
	}

	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return out.Close()
}
