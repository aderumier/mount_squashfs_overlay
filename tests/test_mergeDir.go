package main

import (
	"fmt"
	"strings"
	"path/filepath"
)

func toFSPath(p string) string {
	p = strings.TrimLeft(filepath.ToSlash(p), "/")
	if p == "" {
		return "."
	}
	return p
}

func main() {
	p1 := toFSPath("\\")
	fmt.Printf("toFSPath('\\\\') = %q\n", p1)

	dirPrefix := strings.TrimRight("\\", "/")
	fmt.Printf("dirPrefix = %q\n", dirPrefix)

	childPath := dirPrefix + "/" + "autorun.cmd"
	fmt.Printf("childPath = %q\n", childPath)

	p2 := toFSPath(childPath)
	fmt.Printf("toFSPath(childPath) = %q\n", p2)
}
