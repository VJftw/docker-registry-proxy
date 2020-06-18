package installer

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var (
	ErrNotRegularFile = errors.New("not a regular file")
)

// Copy copies a given source file to the given destination path
func Copy(src string, dest string) error {
	srcStat, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("could not state file '%s': %w", src, err)
	}

	if !srcStat.Mode().IsRegular() {
		return fmt.Errorf("'%s': %w", src, ErrNotRegularFile)
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("could not open source '%s': %w", src, err)
	}
	defer srcFile.Close()

	destDir := filepath.Dir(dest)
	if err := os.MkdirAll(destDir, os.ModeDir); err != nil {
		return fmt.Errorf("could not ensure directory '%s': %w", destDir, err)
	}

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("could not create destination '%s': %w", dest, err)
	}
	defer destFile.Close()

	if _, err = io.Copy(destFile, srcFile); err != nil {
		return fmt.Errorf("could not copy '%s' to '%s': %w", src, dest, err)
	}

	if err := destFile.Chmod(os.ModePerm); err != nil {
		return fmt.Errorf("could not make '%s' executable: %w", dest, err)
	}

	return nil
}
