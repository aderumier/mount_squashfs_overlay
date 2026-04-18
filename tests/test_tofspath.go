package main

import (
	"fmt"
	"strings"
)

func toFSPath(fusePath string) string {
	p := strings.ReplaceAll(fusePath, "\\", "/")
	p = strings.TrimPrefix(p, "/")
	for len(p) > 0 && p[0] == '/' {
		p = p[1:]
	}
	if p == "" {
		return "."
	}
	return p
}

func main() {
	fmt.Println("toFSPath(\"\"):", toFSPath("\\"))
	fmt.Println("toFSPath(\"/\"):", toFSPath("/"))
	fmt.Println("toFSPath(\"\\\\dir\\\\file\"):", toFSPath("\\dir\\file"))
}