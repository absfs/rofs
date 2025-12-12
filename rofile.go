package rofs

import (
	"io/fs"
	"os"

	"github.com/absfs/absfs"
)

type File struct {
	f absfs.File
}

func (f *File) Name() string {
	return f.f.Name()
}

func (f *File) Read(p []byte) (int, error) {
	return f.f.Read(p)
}

func (f *File) ReadAt(b []byte, off int64) (n int, err error) {
	return f.f.ReadAt(b, off)
}

func (f *File) Write(p []byte) (int, error) {
	return 0, &os.PathError{Op: "write", Path: f.f.Name(), Err: os.ErrPermission}
}

func (f *File) WriteAt(b []byte, off int64) (n int, err error) {
	return 0, &os.PathError{Op: "write", Path: f.f.Name(), Err: os.ErrPermission}
}

func (f *File) Close() error {
	return f.f.Close()
}

func (f *File) Seek(offset int64, whence int) (ret int64, err error) {
	return f.f.Seek(offset, whence)
}

func (f *File) Stat() (os.FileInfo, error) {
	return f.f.Stat()
}

func (f *File) Sync() error {
	return nil
}

func (f *File) Readdir(n int) ([]os.FileInfo, error) {
	return f.f.Readdir(n)
}

func (f *File) Readdirnames(n int) ([]string, error) {
	return f.f.Readdirnames(n)
}

func (f *File) Truncate(size int64) error {
	return &os.PathError{Op: "write", Path: f.f.Name(), Err: os.ErrPermission}
}

func (f *File) WriteString(s string) (n int, err error) {
	return 0, &os.PathError{Op: "write", Path: f.f.Name(), Err: os.ErrPermission}
}

// ReadDir reads the contents of the directory and returns a slice of up to n
// DirEntry values. This is a read operation, so it's allowed in read-only mode.
func (f *File) ReadDir(n int) ([]fs.DirEntry, error) {
	return f.f.ReadDir(n)
}
