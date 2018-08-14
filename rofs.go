package rofs

import (
	"os"
	"time"

	"github.com/absfs/absfs"
)

type FileSystem struct {
	fs absfs.FileSystem
}

func NewFS(fs absfs.FileSystem) (*FileSystem, error) {
	return &FileSystem{fs}, nil
}

// FileSystem interface

// OpenFile opens a file using the given flags and the given mode.
func (f *FileSystem) OpenFile(name string, flag int, perm os.FileMode) (absfs.File, error) {
	// error if access mode is not readonly
	if flag&absfs.O_ACCESS != os.O_RDONLY {
		return nil, os.ErrPermission
	}

	file, err := f.fs.OpenFile(name, flag, perm)
	return &File{file}, err
}

// Mkdir creates a directory in the filesystem, return an error if any
// happens.
func (f *FileSystem) Mkdir(name string, perm os.FileMode) error {
	return os.ErrPermission
}

// Remove removes a file identified by name, returning an error, if any
// happens.
func (f *FileSystem) Remove(name string) error {
	return os.ErrPermission
}

// Stat returns the FileInfo structure describing file. If there is an error,
// it will be of type *PathError.
func (f *FileSystem) Stat(name string) (os.FileInfo, error) {
	return f.fs.Stat(name)
}

//Chmod changes the mode of the named file to mode.
func (f *FileSystem) Chmod(name string, mode os.FileMode) error {
	return os.ErrPermission
}

//Chtimes changes the access and modification times of the named file
func (f *FileSystem) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return os.ErrPermission
}

//Chown changes the owner and group ids of the named file
func (f *FileSystem) Chown(name string, uid, gid int) error {
	return os.ErrPermission
}

func (f *FileSystem) Separator() uint8 {
	return f.fs.Separator()
}

func (f *FileSystem) ListSeparator() uint8 {
	return f.fs.ListSeparator()
}

func (f *FileSystem) Chdir(dir string) error {
	return f.fs.Chdir(dir)
}

func (f *FileSystem) Getwd() (dir string, err error) {
	return f.fs.Getwd()
}

func (f *FileSystem) TempDir() string {
	return f.fs.TempDir()
}

func (f *FileSystem) Open(name string) (absfs.File, error) {
	return f.OpenFile(name, os.O_RDONLY, 0444)
}

func (f *FileSystem) Create(name string) (absfs.File, error) {
	return nil, os.ErrPermission
}

func (f *FileSystem) MkdirAll(name string, perm os.FileMode) error {
	return os.ErrPermission
}

func (f *FileSystem) RemoveAll(path string) (err error) {
	return os.ErrPermission
}

func (f *FileSystem) Truncate(name string, size int64) error {
	return os.ErrPermission
}
