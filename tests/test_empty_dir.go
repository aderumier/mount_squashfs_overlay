package main
import (
	"fmt"
	"os"
)
func main() {
	os.Mkdir("empty_test", 0755)
	f, _ := os.Open("empty_test")
	entries, err := f.Readdir(-1)
	fmt.Printf("entries: %d, err: %v\n", len(entries), err)
}
