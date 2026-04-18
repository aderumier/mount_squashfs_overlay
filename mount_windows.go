//go:build windows

package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/winfsp/go-winfsp"
	"github.com/winfsp/go-winfsp/gofs"
)

// Mount mounts the overlay filesystem at the given Windows drive letter (e.g. "Z:").
// This call blocks until the process receives an interrupt signal or is killed.
func Mount(sq *SquashLayer, upperDir string, drive string, debug bool) error {
	var logFile *os.File
	if debug {
		logPath := filepath.Join(filepath.Dir(upperDir), "squashoverlay.log")
		logFile, _ = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	}

	overlayFS := &OverlayFileSystem{
		squash:   sq,
		upperDir: upperDir,
		debug:    debug,
		logFile:  logFile,
	}

	if debug {
		fmt.Fprintf(os.Stderr, "Mounting with go-winfsp at %s ...\n", drive)
	}

	ptfs, err := winfsp.Mount(
		gofs.New(overlayFS),
		drive,
	)
	if err != nil {
		return fmt.Errorf("WinFsp mount failed for %s: %w", drive, err)
	}
	defer ptfs.Unmount()

	if debug {
		fmt.Fprintf(os.Stderr, "Mounted successfully. Waiting for interrupt...\n")
	}

	// Block until the process is killed or interrupted.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	if debug {
		fmt.Fprintf(os.Stderr, "Unmounting...\n")
	}
	return nil
}

// Umount is a no-op on the go-winfsp path; the mount is cleaned up
// when the process exits. Provided for CLI interface compatibility.
func Umount(drive string) error {
	return nil
}
