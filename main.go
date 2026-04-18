package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const version = "0.1.0"

// CLI interface is designed as a drop-in replacement for the mount.exe used by
// EmulatorLauncher (github.com/RetroBat-Official/emulatorlauncher).
//
// Invocation by the launcher:
//
//	mount.exe [-debug] -drive Z: [-extractionpath "<path>"] [-overlay "<path>"] "<squashfs-file>"
//
// The process runs until killed; killing it unmounts the drive.
func main() {
	debug := flag.Bool("debug", false, "enable verbose debug output")
	drive := flag.String("drive", "", "drive letter to mount at, e.g. Z:")
	extractionPath := flag.String("extractionpath", "", "extraction/work directory")
	overlayPath := flag.String("overlay", "", "persistent writable overlay directory")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "error: squashfs file argument is required\n\n")
		usage()
		os.Exit(1)
	}
	if *drive == "" {
		fmt.Fprintf(os.Stderr, "error: -drive is required\n\n")
		usage()
		os.Exit(1)
	}

	squashFile, err := filepath.Abs(flag.Arg(0))
	if err != nil {
		fatalf("invalid squashfs path: %v", err)
	}
	if _, err := os.Stat(squashFile); err != nil {
		fatalf("cannot open squashfs file %q: %v", squashFile, err)
	}

	driveLetter := normalizeDrive(*drive)

	// -overlay is the writable upper directory; -extractionpath is accepted
	// for compatibility with existing EmulatorLauncher invocations but ignored.
	upperDir := *overlayPath
	if upperDir == "" {
		// Last resort: place overlay next to the squashfs file.
		upperDir = squashFile + ".overlay"
	}
	upperDir = filepath.Join(upperDir, "upper")
	if err := os.MkdirAll(upperDir, 0755); err != nil {
		fatalf("cannot create overlay dir %q: %v", upperDir, err)
	}

	if *debug {
		fmt.Fprintf(os.Stderr, "squashoverlay v%s\n", version)
		fmt.Fprintf(os.Stderr, "  squashfs    : %s\n", squashFile)
		fmt.Fprintf(os.Stderr, "  drive       : %s\n", driveLetter)
		fmt.Fprintf(os.Stderr, "  upper dir   : %s\n", upperDir)
		fmt.Fprintf(os.Stderr, "  extractpath : %s\n", *extractionPath)
		fmt.Fprintf(os.Stderr, "  overlay arg : %s\n", *overlayPath)
	}

	if err := checkWinFsp(); err != nil {
		fatalf("WinFsp not available: %v\n\nDownload and install WinFsp from:\nhttps://github.com/winfsp/winfsp/releases", err)
	}

	sq, err := NewSquashLayer(squashFile)
	if err != nil {
		fatalf("failed to open squashfs %q: %v", squashFile, err)
	}

	if *debug {
		fmt.Fprintf(os.Stderr, "Mounting %s at %s ...\n", squashFile, driveLetter)
	}

	// Mount blocks until the filesystem is unmounted (i.e. this process is killed).
	if err := Mount(sq, upperDir, driveLetter, *debug); err != nil {
		fatalf("mount failed: %v", err)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `squashoverlay v%s - mount a squashfs file as a Windows drive with persistent writable overlay

Usage:
  squashoverlay.exe [-debug] -drive <X:> [-extractionpath <dir>] [-overlay <dir>] <squashfs-file>

Flags:
  -drive <X:>            Drive letter to mount at (required)
  -extractionpath <dir>  Work/extraction directory (used as overlay if -overlay not given)
  -overlay <dir>         Persistent writable overlay directory (takes precedence)
  -debug                 Verbose output

The process runs until killed; killing it unmounts the drive.
Requires WinFsp >= 1.10: https://github.com/winfsp/winfsp/releases
`, version)
}

func normalizeDrive(s string) string {
	s = strings.TrimSuffix(s, `\`)
	s = strings.TrimSuffix(s, "/")
	s = strings.ToUpper(s)
	if !strings.HasSuffix(s, ":") {
		s += ":"
	}
	return s
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}
