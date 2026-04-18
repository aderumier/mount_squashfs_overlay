package main

import (
	"fmt"
	"os"
	"time"

	"github.com/winfsp/go-winfsp"
	"github.com/winfsp/go-winfsp/gofs"
)

type MemFS struct{}

func (f *MemFS) Stat(name string) (os.FileInfo, error) {
	fmt.Printf("Stat: %q\n", name)
	if name == "" || name == "\\" || name == "/" {
		return &MemFileInfo{name: "root", isDir: true}, nil
	}
	return nil, os.ErrNotExist
}
func (f *MemFS) OpenFile(name string, flag int, perm os.FileMode) (gofs.File, error) {
	fmt.Printf("OpenFile: %q\n", name)
	if name == "" || name == "\\" || name == "/" {
		return &MemFile{path: name, isDir: true}, nil
	}
	return nil, os.ErrNotExist
}
func (f *MemFS) Mkdir(name string, perm os.FileMode) error { return os.ErrPermission }
func (f *MemFS) Remove(name string) error                  { return os.ErrPermission }
func (f *MemFS) Rename(source string, target string) error { return os.ErrPermission }

type MemFile struct {
	path  string
	isDir bool
}

func (f *MemFile) Read(p []byte) (int, error)               { return 0, os.ErrPermission }
func (f *MemFile) ReadAt(p []byte, off int64) (int, error)  { return 0, os.ErrPermission }
func (f *MemFile) Write(p []byte) (int, error)              { return 0, os.ErrPermission }
func (f *MemFile) WriteAt(p []byte, off int64) (int, error) { return 0, os.ErrPermission }
func (f *MemFile) Seek(offset int64, whence int) (int64, error) { return 0, os.ErrPermission }
func (f *MemFile) Close() error                             { return nil }
func (f *MemFile) Truncate(size int64) error                { return os.ErrPermission }
func (f *MemFile) Sync() error                              { return nil }

func (f *MemFile) Stat() (os.FileInfo, error) {
	fmt.Printf("File.Stat: %q\n", f.path)
	return &MemFileInfo{name: f.path, isDir: f.isDir}, nil
}

func (f *MemFile) Readdir(n int) ([]os.FileInfo, error) {
	fmt.Printf("File.Readdir: %q\n", f.path)
	return []os.FileInfo{
		&MemFileInfo{name: "test.txt", isDir: false},
	}, nil
}

type MemFileInfo struct {
	name  string
	isDir bool
}

func (i *MemFileInfo) Name() string       { return i.name }
func (i *MemFileInfo) Size() int64        { return 100 }
func (i *MemFileInfo) Mode() os.FileMode  {
	if i.isDir {
		return os.ModeDir | 0755
	}
	return 0644
}
func (i *MemFileInfo) ModTime() time.Time { return time.Now() }
func (i *MemFileInfo) IsDir() bool        { return i.isDir }
func (i *MemFileInfo) Sys() any           { return nil }

func main() {
	fsys := &MemFS{}
	ptfs, err := winfsp.Mount(gofs.New(fsys), "Y:")
	if err != nil {
		panic(err)
	}
	defer ptfs.Unmount()
	select {}
}
