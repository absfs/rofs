package ptfs

import (
	"os"
	"time"

	"github.com/absfs/absfs"
)

type Filer struct {
	fs absfs.Filer
}

func NewFiler(fs absfs.Filer) (*Filer, error) {
	return &Filer{fs}, nil
}

// Filer interface

// OpenFile opens a file using the given flags and the given mode.
func (f *Filer) OpenFile(name string, flag int, perm os.FileMode) (absfs.File, error) {
	return f.fs.OpenFile(name, flag, perm)
}

// Mkdir creates a directory in the filesystem, return an error if any
// happens.
func (f *Filer) Mkdir(name string, perm os.FileMode) error {
	return f.fs.Mkdir(name, perm)
}

// Remove removes a file identified by name, returning an error, if any
// happens.
func (f *Filer) Remove(name string) error {
	return f.fs.Remove(name)
}

// Stat returns the FileInfo structure describing file. If there is an error,
// it will be of type *PathError.
func (f *Filer) Stat(name string) (os.FileInfo, error) {
	return f.fs.Stat(name)
}

//Chmod changes the mode of the named file to mode.
func (f *Filer) Chmod(name string, mode os.FileMode) error {
	return f.fs.Chmod(name, mode)
}

//Chtimes changes the access and modification times of the named file
func (f *Filer) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return f.fs.Chtimes(name, atime, mtime)
}

//Chown changes the owner and group ids of the named file
func (f *Filer) Chown(name string, uid, gid int) error {
	return f.fs.Chown(name, uid, gid)
}

type FileSystem struct {
	fs absfs.FileSystem
}

func NewFS(fs absfs.FileSystem) (*FileSystem, error) {
	return &FileSystem{fs}, nil
}

// FileSystem interface

// OpenFile opens a file using the given flags and the given mode.
func (f *FileSystem) OpenFile(name string, flag int, perm os.FileMode) (absfs.File, error) {
	flag = os.O_RDONLY | flag&^absfs.O_ACCESS
	if flag&absfs.O_ACCESS != os.O_RDONLY {
		panic("fail")
	}
	return f.fs.OpenFile(name, flag, perm)
}

// Mkdir creates a directory in the filesystem, return an error if any
// happens.
func (f *FileSystem) Mkdir(name string, perm os.FileMode) error {
	return f.fs.Mkdir(name, perm)
}

// Remove removes a file identified by name, returning an error, if any
// happens.
func (f *FileSystem) Remove(name string) error {
	return f.fs.Remove(name)
}

// Stat returns the FileInfo structure describing file. If there is an error,
// it will be of type *PathError.
func (f *FileSystem) Stat(name string) (os.FileInfo, error) {
	return f.fs.Stat(name)
}

//Chmod changes the mode of the named file to mode.
func (f *FileSystem) Chmod(name string, mode os.FileMode) error {
	return f.fs.Chmod(name, mode)
}

//Chtimes changes the access and modification times of the named file
func (f *FileSystem) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return f.fs.Chtimes(name, atime, mtime)
}

//Chown changes the owner and group ids of the named file
func (f *FileSystem) Chown(name string, uid, gid int) error {
	return f.fs.Chown(name, uid, gid)
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
	return f.fs.Open(name)
}

func (f *FileSystem) Create(name string) (absfs.File, error) {
	return f.fs.Create(name)
}

func (f *FileSystem) MkdirAll(name string, perm os.FileMode) error {
	return f.fs.MkdirAll(name, perm)
}

func (f *FileSystem) RemoveAll(path string) (err error) {
	return f.fs.RemoveAll(path)
}

func (f *FileSystem) Truncate(name string, size int64) error {
	return f.fs.Truncate(name, size)
}

type SymlinkFileSystem struct {
	sfs absfs.SymlinkFileSystem
}

func NewSymlinkFS(fs absfs.SymlinkFileSystem) (*SymlinkFileSystem, error) {
	return &SymlinkFileSystem{fs}, nil
}

// OpenFile opens a file using the given flags and the given mode.
func (f *SymlinkFileSystem) OpenFile(name string, flag int, perm os.FileMode) (absfs.File, error) {
	return f.sfs.OpenFile(name, flag, perm)
}

// Mkdir creates a directory in the filesystem, return an error if any
// happens.
func (f *SymlinkFileSystem) Mkdir(name string, perm os.FileMode) error {
	return f.sfs.Mkdir(name, perm)
}

// Remove removes a file identified by name, returning an error, if any
// happens.
func (f *SymlinkFileSystem) Remove(name string) error {
	return f.sfs.Remove(name)
}

// Stat returns the FileInfo structure describing file. If there is an error,
// it will be of type *PathError.
func (f *SymlinkFileSystem) Stat(name string) (os.FileInfo, error) {
	return f.sfs.Stat(name)
}

//Chmod changes the mode of the named file to mode.
func (f *SymlinkFileSystem) Chmod(name string, mode os.FileMode) error {
	return f.sfs.Chmod(name, mode)
}

//Chtimes changes the access and modification times of the named file
func (f *SymlinkFileSystem) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return f.sfs.Chtimes(name, atime, mtime)
}

//Chown changes the owner and group ids of the named file
func (f *SymlinkFileSystem) Chown(name string, uid, gid int) error {
	return f.sfs.Chown(name, uid, gid)
}

func (f *SymlinkFileSystem) Separator() uint8 {
	return f.sfs.Separator()
}

func (f *SymlinkFileSystem) ListSeparator() uint8 {
	return f.sfs.ListSeparator()
}

func (f *SymlinkFileSystem) Chdir(dir string) error {
	return f.sfs.Chdir(dir)
}

func (f *SymlinkFileSystem) Getwd() (dir string, err error) {
	return f.sfs.Getwd()
}

func (f *SymlinkFileSystem) TempDir() string {
	return f.sfs.TempDir()
}

func (f *SymlinkFileSystem) Open(name string) (absfs.File, error) {
	return f.sfs.Open(name)
}

func (f *SymlinkFileSystem) Create(name string) (absfs.File, error) {
	return f.sfs.Create(name)
}

func (f *SymlinkFileSystem) MkdirAll(name string, perm os.FileMode) error {
	return f.sfs.MkdirAll(name, perm)
}

func (f *SymlinkFileSystem) RemoveAll(path string) (err error) {
	return f.sfs.RemoveAll(path)
}

func (f *SymlinkFileSystem) Truncate(name string, size int64) error {
	return f.sfs.Truncate(name, size)
}

// Lstat returns a FileInfo describing the named file. If the file is a
// symbolic link, the returned FileInfo describes the symbolic link. Lstat
// makes no attempt to follow the link. If there is an error, it will be of type *PathError.
func (f *SymlinkFileSystem) Lstat(name string) (os.FileInfo, error) {
	return f.sfs.Lstat(name)
}

// Lchown changes the numeric uid and gid of the named file. If the file is a
// symbolic link, it changes the uid and gid of the link itself. If there is
// an error, it will be of type *PathError.
//
// On Windows, it always returns the syscall.EWINDOWS error, wrapped in
// *PathError.
func (f *SymlinkFileSystem) Lchown(name string, uid, gid int) error {
	return f.sfs.Lchown(name, uid, gid)
}

// Readlink returns the destination of the named symbolic link. If there is an
// error, it will be of type *PathError.
func (f *SymlinkFileSystem) Readlink(name string) (string, error) {
	return f.sfs.Readlink(name)
}

// Symlink creates newname as a symbolic link to oldname. If there is an
// error, it will be of type *LinkError.
func (f *SymlinkFileSystem) Symlink(oldname, newname string) error {
	return f.sfs.Symlink(oldname, newname)
}
