package main
import (
	"fmt"
	"os"
)
func main() {
	sq, err := NewSquashLayer(os.Args[1])
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	os.MkdirAll("test_upper/upper", 0755)
	
	ofs := &OverlayFileSystem{
		squash:   sq,
		upperDir: "test_upper/upper",
	}
	
	// manually emulate what go-winfsp does:
	f, err := ofs.OpenFile("", os.O_RDONLY, 0)
	if err != nil {
		fmt.Printf("OpenFile err: %v\n", err)
		return
	}
	
	entries, err := f.Readdir(-1)
	if err != nil {
		fmt.Printf("Readdir err: %v\n", err)
		return
	}
	
	fmt.Printf("overlayFile.Readdir(-1) returned %d entries:\n", len(entries))
	for _, e := range entries {
		fmt.Printf("  %s\n", e.Name())
	}
}
