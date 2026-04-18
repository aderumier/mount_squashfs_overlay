//go:build windows
package main
import (
"fmt"
"os"
"path/filepath"
"github.com/winfsp/go-winfsp"
"github.com/winfsp/go-winfsp/gofs"
)
type MockFS struct {}
func (m *MockFS) OpenFile(name string, flag int, perm os.FileMode) (gofs.File, error) {
fmt.Fprintf(os.Stderr, "OpenFile called on %q\n", name)
return nil, os.ErrNotExist
}
func (m *MockFS) Mkdir(name string, perm os.FileMode) error { return nil }
func (m *MockFS) Stat(name string) (os.FileInfo, error) {
fmt.Fprintf(os.Stderr, "Stat called on %q\n", name)
return nil, os.ErrNotExist
}
func (m *MockFS) Rename(s, t string) error { return nil }
func (m *MockFS) Remove(name string) error { return nil }

func main() {
winfsp.Mount(gofs.New(&MockFS{}), "Z:")
}
