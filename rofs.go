package rofs

import (
	"io/fs"
	"os"
	"time"

	"github.com/absfs/absfs"
)

type FileSystem struct {
	fs absfs.SymlinkFileSystem
}

func NewFS(fs absfs.SymlinkFileSystem) (*FileSystem, error) {
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

// Rename renames (moves) oldpath to newpath. If newpath already exists and
// is not a directory, Rename replaces it. OS-specific restrictions may apply
// when oldpath and newpath are in different directories. If there is an
// error, it will be of type *LinkError.
func (f *FileSystem) Rename(oldpath, newpath string) error {
	return &os.LinkError{
		Op:  "rename",
		Old: oldpath,
		New: newpath,
		Err: os.ErrPermission,
	}
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

// Lstat returns a FileInfo describing the named file. If the file is a
// symbolic link, the returned FileInfo describes the symbolic link. Lstat
// makes no attempt to follow the link. If there is an error, it will be of type *PathError.
func (f *FileSystem) Lstat(name string) (os.FileInfo, error) {
	return f.fs.Lstat(name)
}

// Lchown changes the numeric uid and gid of the named file. If the file is a
// symbolic link, it changes the uid and gid of the link itself. If there is
// an error, it will be of type *PathError.
//
// On Windows, it always returns the syscall.EWINDOWS error, wrapped in
// `*PathError`.
func (f *FileSystem) Lchown(name string, uid, gid int) error {
	// return f.fs.Lchown(name, uid, gid)
	return &os.PathError{
		Op:   "lchown",
		Path: name,
		Err:  os.ErrPermission,
	}
}

// Readlink returns the destination of the named symbolic link. If there is an
// error, it will be of type *PathError.
func (f *FileSystem) Readlink(name string) (string, error) {
	return f.fs.Readlink(name)
}

// Symlink creates newname as a symbolic link to oldname. If there is an
// error, it will be of type *LinkError.
func (f *FileSystem) Symlink(oldname, newname string) error {
	return &os.LinkError{
		Op:  "symlink",
		Old: oldname,
		New: newname,
		Err: os.ErrPermission,
	}
}

// ReadDir reads the named directory and returns a list of directory entries.
// This is a read operation, so it's allowed in read-only mode.
func (f *FileSystem) ReadDir(name string) ([]fs.DirEntry, error) {
	return f.fs.ReadDir(name)
}

// ReadFile reads the named file and returns its contents.
// This is a read operation, so it's allowed in read-only mode.
func (f *FileSystem) ReadFile(name string) ([]byte, error) {
	return f.fs.ReadFile(name)
}

// Sub returns an fs.FS corresponding to the subtree rooted at dir.
// The result is wrapped in rofs to maintain read-only guarantee.
func (f *FileSystem) Sub(dir string) (fs.FS, error) {
	return absfs.FilerToFS(f, dir)
}
