package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// migrateDeletions converts a RetroBat-style .deletions text file to
// Docker/OCI whiteout files (.wh.<name>) and removes the original file.
//
// The original format stores one deleted path per line using backslash
// separators, e.g. "\\subdir\\file.rom". Each entry is converted to a
// zero-byte whiteout file at <overlayDir>/<parent>/.wh.<name>.
func migrateDeletions(overlayDir string) error {
	deletionsFile := filepath.Join(overlayDir, ".deletions")
	f, err := os.Open(deletionsFile)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		// Normalize: replace backslashes, strip leading slashes.
		p := strings.ReplaceAll(line, "\\", "/")
		p = strings.TrimLeft(p, "/")
		if p == "" {
			continue
		}

		dir := filepath.Dir(p)
		base := filepath.Base(p)
		whiteout := filepath.Join(overlayDir, filepath.FromSlash(dir), ".wh."+base)

		if err := os.MkdirAll(filepath.Dir(whiteout), 0755); err != nil {
			return err
		}
		wf, err := os.Create(whiteout)
		if err != nil {
			return err
		}
		wf.Close()
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	f.Close()
	return os.Remove(deletionsFile)
}
