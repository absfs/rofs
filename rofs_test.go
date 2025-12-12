package rofs_test

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"
	"time"

	"github.com/absfs/absfs"
	"github.com/absfs/ioutil"
	"github.com/absfs/memfs"
	"github.com/absfs/osfs"
	"github.com/absfs/rofs"
)

func TestInterface(t *testing.T) {
	osfs, err := osfs.NewFS()
	if err != nil {
		t.Fatal(err)
	}
	testfs, err := rofs.NewFS(osfs)
	if err != nil {
		t.Fatal(err)
	}

	var fs absfs.SymlinkFileSystem
	fs = testfs
	_ = fs
}

func TestPtfs(t *testing.T) {
	wfs, err := memfs.NewFS()
	if err != nil {
		t.Fatal(err)
	}

	err = wfs.Mkdir("/write_here", 0777)
	if err != nil {
		t.Fatal(err)
	}

	dataOut := []byte("Bet you can't change me!")
	err = ioutil.WriteFile(wfs, "/write_here/file.txt", dataOut, 0666)
	if err != nil {
		t.Fatal(err)
	}

	rfs, err := rofs.NewFS(wfs)
	if err != nil {
		t.Fatal(err)
	}

	err = rfs.Mkdir("/write_here", 0777)
	if err == nil {
		t.Fatal("expected errors, but got go none")
	}

	err = rfs.Mkdir("/can_I_make_a_dir", 0777)
	if err == nil {
		root, err := wfs.Open("/")
		if err != nil {
			t.Fatal(err)
		}
		defer root.Close()
		list, err := root.Readdir(-1)
		if err != nil {
			t.Fatal(err)
		}

		for i, info := range list {
			t.Logf("%d: %s", i, info.Name())
		}
		t.Fatal("expected errors, but got go none")
	}

	// I should be able to read a file.
	dataIn, err := ioutil.ReadFile(rfs, "/write_here/file.txt")
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(dataOut, dataIn) {
		t.Fatalf("expected read to match write: %q expect %q", string(dataIn), string(dataOut))
	}

	dataIn = []byte("I can change you!")
	err = ioutil.WriteFile(rfs, "/write_here/file.txt", dataIn, 0666)
	if err == nil {
		t.Fatal("expected errors, but got go none")
	}

	err = rfs.Remove("/write_here/file.txt")
	if err == nil {
		t.Fatal("expected errors, but got go none")
	}

	info, err := rfs.Stat("/write_here/file.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s %s % 6d %s", info.Mode(), info.Name(), info.Size(), info.ModTime().Format(time.UnixDate))
}

// setupTestFS creates a memfs with test files and wraps it with rofs
func setupTestFS(t *testing.T) (*rofs.FileSystem, absfs.SymlinkFileSystem) {
	t.Helper()
	wfs, err := memfs.NewFS()
	if err != nil {
		t.Fatal(err)
	}

	// Create test directories
	if err := wfs.MkdirAll("/testdir/subdir", 0755); err != nil {
		t.Fatal(err)
	}

	// Create test files
	if err := ioutil.WriteFile(wfs, "/testdir/file.txt", []byte("test content"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(wfs, "/testdir/subdir/nested.txt", []byte("nested content"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(wfs, "/empty.txt", []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	// Create symlink
	if err := wfs.Symlink("/testdir/file.txt", "/testdir/link.txt"); err != nil {
		t.Fatal(err)
	}

	rfs, err := rofs.NewFS(wfs)
	if err != nil {
		t.Fatal(err)
	}

	return rfs, wfs
}

// Phase 2: Test FileSystem Write Operations Return Errors

func TestFileSystemWriteOperationsReturnPermissionError(t *testing.T) {
	rfs, _ := setupTestFS(t)

	t.Run("Mkdir returns ErrPermission", func(t *testing.T) {
		err := rfs.Mkdir("/newdir", 0755)
		if !errors.Is(err, os.ErrPermission) {
			t.Errorf("Mkdir: expected os.ErrPermission, got %v", err)
		}
	})

	t.Run("Remove returns ErrPermission", func(t *testing.T) {
		err := rfs.Remove("/testdir/file.txt")
		if !errors.Is(err, os.ErrPermission) {
			t.Errorf("Remove: expected os.ErrPermission, got %v", err)
		}
	})

	t.Run("Rename returns LinkError with ErrPermission", func(t *testing.T) {
		err := rfs.Rename("/testdir/file.txt", "/testdir/renamed.txt")
		var linkErr *os.LinkError
		if !errors.As(err, &linkErr) {
			t.Errorf("Rename: expected *os.LinkError, got %T", err)
		} else {
			if linkErr.Op != "rename" {
				t.Errorf("LinkError.Op: expected 'rename', got %q", linkErr.Op)
			}
			if linkErr.Old != "/testdir/file.txt" {
				t.Errorf("LinkError.Old: expected '/testdir/file.txt', got %q", linkErr.Old)
			}
			if linkErr.New != "/testdir/renamed.txt" {
				t.Errorf("LinkError.New: expected '/testdir/renamed.txt', got %q", linkErr.New)
			}
			if !errors.Is(linkErr.Err, os.ErrPermission) {
				t.Errorf("LinkError.Err: expected os.ErrPermission, got %v", linkErr.Err)
			}
		}
	})

	t.Run("Chmod returns ErrPermission", func(t *testing.T) {
		err := rfs.Chmod("/testdir/file.txt", 0644)
		if !errors.Is(err, os.ErrPermission) {
			t.Errorf("Chmod: expected os.ErrPermission, got %v", err)
		}
	})

	t.Run("Chtimes returns ErrPermission", func(t *testing.T) {
		now := time.Now()
		err := rfs.Chtimes("/testdir/file.txt", now, now)
		if !errors.Is(err, os.ErrPermission) {
			t.Errorf("Chtimes: expected os.ErrPermission, got %v", err)
		}
	})

	t.Run("Chown returns ErrPermission", func(t *testing.T) {
		err := rfs.Chown("/testdir/file.txt", 1000, 1000)
		if !errors.Is(err, os.ErrPermission) {
			t.Errorf("Chown: expected os.ErrPermission, got %v", err)
		}
	})

	t.Run("Create returns ErrPermission", func(t *testing.T) {
		f, err := rfs.Create("/newfile.txt")
		if f != nil {
			t.Error("Create: expected nil file")
		}
		if !errors.Is(err, os.ErrPermission) {
			t.Errorf("Create: expected os.ErrPermission, got %v", err)
		}
	})

	t.Run("MkdirAll returns ErrPermission", func(t *testing.T) {
		err := rfs.MkdirAll("/new/nested/dir", 0755)
		if !errors.Is(err, os.ErrPermission) {
			t.Errorf("MkdirAll: expected os.ErrPermission, got %v", err)
		}
	})

	t.Run("RemoveAll returns ErrPermission", func(t *testing.T) {
		err := rfs.RemoveAll("/testdir")
		if !errors.Is(err, os.ErrPermission) {
			t.Errorf("RemoveAll: expected os.ErrPermission, got %v", err)
		}
	})

	t.Run("Truncate returns ErrPermission", func(t *testing.T) {
		err := rfs.Truncate("/testdir/file.txt", 0)
		if !errors.Is(err, os.ErrPermission) {
			t.Errorf("Truncate: expected os.ErrPermission, got %v", err)
		}
	})

	t.Run("Lchown returns PathError with ErrPermission", func(t *testing.T) {
		err := rfs.Lchown("/testdir/link.txt", 1000, 1000)
		var pathErr *os.PathError
		if !errors.As(err, &pathErr) {
			t.Errorf("Lchown: expected *os.PathError, got %T", err)
		} else {
			if pathErr.Op != "lchown" {
				t.Errorf("PathError.Op: expected 'lchown', got %q", pathErr.Op)
			}
			if pathErr.Path != "/testdir/link.txt" {
				t.Errorf("PathError.Path: expected '/testdir/link.txt', got %q", pathErr.Path)
			}
			if !errors.Is(pathErr.Err, os.ErrPermission) {
				t.Errorf("PathError.Err: expected os.ErrPermission, got %v", pathErr.Err)
			}
		}
	})

	t.Run("Symlink returns LinkError with ErrPermission", func(t *testing.T) {
		err := rfs.Symlink("/testdir/file.txt", "/testdir/newlink.txt")
		var linkErr *os.LinkError
		if !errors.As(err, &linkErr) {
			t.Errorf("Symlink: expected *os.LinkError, got %T", err)
		} else {
			if linkErr.Op != "symlink" {
				t.Errorf("LinkError.Op: expected 'symlink', got %q", linkErr.Op)
			}
			if linkErr.Old != "/testdir/file.txt" {
				t.Errorf("LinkError.Old: expected '/testdir/file.txt', got %q", linkErr.Old)
			}
			if linkErr.New != "/testdir/newlink.txt" {
				t.Errorf("LinkError.New: expected '/testdir/newlink.txt', got %q", linkErr.New)
			}
			if !errors.Is(linkErr.Err, os.ErrPermission) {
				t.Errorf("LinkError.Err: expected os.ErrPermission, got %v", linkErr.Err)
			}
		}
	})
}

// Phase 2: Test File Write Operations Return Errors

func TestFileWriteOperationsReturnPermissionError(t *testing.T) {
	rfs, _ := setupTestFS(t)

	file, err := rfs.Open("/testdir/file.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	t.Run("Write returns PathError with ErrPermission", func(t *testing.T) {
		n, err := file.Write([]byte("new content"))
		if n != 0 {
			t.Errorf("Write: expected 0 bytes, got %d", n)
		}
		var pathErr *os.PathError
		if !errors.As(err, &pathErr) {
			t.Errorf("Write: expected *os.PathError, got %T", err)
		} else {
			if pathErr.Op != "write" {
				t.Errorf("PathError.Op: expected 'write', got %q", pathErr.Op)
			}
			if !errors.Is(pathErr.Err, os.ErrPermission) {
				t.Errorf("PathError.Err: expected os.ErrPermission, got %v", pathErr.Err)
			}
		}
	})

	t.Run("WriteAt returns PathError with ErrPermission", func(t *testing.T) {
		n, err := file.WriteAt([]byte("new content"), 0)
		if n != 0 {
			t.Errorf("WriteAt: expected 0 bytes, got %d", n)
		}
		var pathErr *os.PathError
		if !errors.As(err, &pathErr) {
			t.Errorf("WriteAt: expected *os.PathError, got %T", err)
		} else {
			if pathErr.Op != "write" {
				t.Errorf("PathError.Op: expected 'write', got %q", pathErr.Op)
			}
			if !errors.Is(pathErr.Err, os.ErrPermission) {
				t.Errorf("PathError.Err: expected os.ErrPermission, got %v", pathErr.Err)
			}
		}
	})

	t.Run("WriteString returns PathError with ErrPermission", func(t *testing.T) {
		n, err := file.WriteString("new content")
		if n != 0 {
			t.Errorf("WriteString: expected 0 bytes, got %d", n)
		}
		var pathErr *os.PathError
		if !errors.As(err, &pathErr) {
			t.Errorf("WriteString: expected *os.PathError, got %T", err)
		} else {
			if pathErr.Op != "write" {
				t.Errorf("PathError.Op: expected 'write', got %q", pathErr.Op)
			}
			if !errors.Is(pathErr.Err, os.ErrPermission) {
				t.Errorf("PathError.Err: expected os.ErrPermission, got %v", pathErr.Err)
			}
		}
	})

	t.Run("Truncate returns PathError with ErrPermission", func(t *testing.T) {
		err := file.Truncate(0)
		var pathErr *os.PathError
		if !errors.As(err, &pathErr) {
			t.Errorf("Truncate: expected *os.PathError, got %T", err)
		} else {
			if pathErr.Op != "write" {
				t.Errorf("PathError.Op: expected 'write', got %q", pathErr.Op)
			}
			if !errors.Is(pathErr.Err, os.ErrPermission) {
				t.Errorf("PathError.Err: expected os.ErrPermission, got %v", pathErr.Err)
			}
		}
	})
}

// Phase 2: Test OpenFile Flags Validation

func TestOpenFileFlagsValidation(t *testing.T) {
	rfs, _ := setupTestFS(t)

	t.Run("O_WRONLY returns ErrPermission", func(t *testing.T) {
		f, err := rfs.OpenFile("/testdir/file.txt", os.O_WRONLY, 0644)
		if f != nil {
			t.Error("expected nil file")
		}
		if !errors.Is(err, os.ErrPermission) {
			t.Errorf("expected os.ErrPermission, got %v", err)
		}
	})

	t.Run("O_RDWR returns ErrPermission", func(t *testing.T) {
		f, err := rfs.OpenFile("/testdir/file.txt", os.O_RDWR, 0644)
		if f != nil {
			t.Error("expected nil file")
		}
		if !errors.Is(err, os.ErrPermission) {
			t.Errorf("expected os.ErrPermission, got %v", err)
		}
	})

	// Note: The rofs implementation only checks the access mode (O_RDONLY, O_WRONLY, O_RDWR),
	// not additional flags like O_APPEND, O_CREATE, O_TRUNC. Those flags are passed through
	// to the underlying filesystem which handles them appropriately.
	t.Run("O_APPEND with O_RDONLY succeeds (access mode is RDONLY)", func(t *testing.T) {
		f, err := rfs.OpenFile("/testdir/file.txt", os.O_RDONLY|os.O_APPEND, 0644)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if f != nil {
			f.Close()
		}
	})

	t.Run("O_WRONLY with O_CREATE returns ErrPermission", func(t *testing.T) {
		f, err := rfs.OpenFile("/newfile.txt", os.O_WRONLY|os.O_CREATE, 0644)
		if f != nil {
			t.Error("expected nil file")
		}
		if !errors.Is(err, os.ErrPermission) {
			t.Errorf("expected os.ErrPermission, got %v", err)
		}
	})

	t.Run("O_WRONLY with O_TRUNC returns ErrPermission", func(t *testing.T) {
		f, err := rfs.OpenFile("/testdir/file.txt", os.O_WRONLY|os.O_TRUNC, 0644)
		if f != nil {
			t.Error("expected nil file")
		}
		if !errors.Is(err, os.ErrPermission) {
			t.Errorf("expected os.ErrPermission, got %v", err)
		}
	})

	t.Run("O_RDONLY succeeds", func(t *testing.T) {
		f, err := rfs.OpenFile("/testdir/file.txt", os.O_RDONLY, 0644)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if f == nil {
			t.Error("expected non-nil file")
		} else {
			f.Close()
		}
	})

	t.Run("Combined write access flags return ErrPermission", func(t *testing.T) {
		// Only test flags where access mode is O_WRONLY or O_RDWR
		flags := []int{
			os.O_WRONLY | os.O_CREATE,
			os.O_RDWR | os.O_APPEND,
			os.O_WRONLY | os.O_TRUNC,
		}
		for _, flag := range flags {
			f, err := rfs.OpenFile("/testdir/file.txt", flag, 0644)
			if f != nil {
				t.Errorf("flag %d: expected nil file", flag)
			}
			if !errors.Is(err, os.ErrPermission) {
				t.Errorf("flag %d: expected os.ErrPermission, got %v", flag, err)
			}
		}
	})
}

// Phase 3: Test File Read Operations

func TestFileReadOperations(t *testing.T) {
	rfs, _ := setupTestFS(t)

	t.Run("Name returns correct filename", func(t *testing.T) {
		file, err := rfs.Open("/testdir/file.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		name := file.Name()
		if name != "/testdir/file.txt" {
			t.Errorf("Name: expected '/testdir/file.txt', got %q", name)
		}
	})

	t.Run("Read reads file content correctly", func(t *testing.T) {
		file, err := rfs.Open("/testdir/file.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		buf := make([]byte, 100)
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			t.Errorf("Read: unexpected error %v", err)
		}
		if !bytes.Equal(buf[:n], []byte("test content")) {
			t.Errorf("Read: expected 'test content', got %q", string(buf[:n]))
		}
	})

	t.Run("ReadAt reads at specific offset", func(t *testing.T) {
		file, err := rfs.Open("/testdir/file.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		buf := make([]byte, 7)
		n, err := file.ReadAt(buf, 5) // "content"
		if err != nil && err != io.EOF {
			t.Errorf("ReadAt: unexpected error %v", err)
		}
		if !bytes.Equal(buf[:n], []byte("content")) {
			t.Errorf("ReadAt: expected 'content', got %q", string(buf[:n]))
		}
	})

	t.Run("Seek changes file position", func(t *testing.T) {
		file, err := rfs.Open("/testdir/file.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		// Seek to position 5
		pos, err := file.Seek(5, io.SeekStart)
		if err != nil {
			t.Errorf("Seek: unexpected error %v", err)
		}
		if pos != 5 {
			t.Errorf("Seek: expected position 5, got %d", pos)
		}

		// Read from position 5
		buf := make([]byte, 7)
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			t.Errorf("Read after Seek: unexpected error %v", err)
		}
		if !bytes.Equal(buf[:n], []byte("content")) {
			t.Errorf("Read after Seek: expected 'content', got %q", string(buf[:n]))
		}
	})

	t.Run("Stat returns correct FileInfo", func(t *testing.T) {
		file, err := rfs.Open("/testdir/file.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		info, err := file.Stat()
		if err != nil {
			t.Errorf("Stat: unexpected error %v", err)
		}
		if info.Name() != "file.txt" {
			t.Errorf("Stat: expected name 'file.txt', got %q", info.Name())
		}
		if info.Size() != int64(len("test content")) {
			t.Errorf("Stat: expected size %d, got %d", len("test content"), info.Size())
		}
		if info.IsDir() {
			t.Error("Stat: expected file, got directory")
		}
	})

	t.Run("Readdir lists directory contents", func(t *testing.T) {
		dir, err := rfs.Open("/testdir")
		if err != nil {
			t.Fatal(err)
		}
		defer dir.Close()

		entries, err := dir.Readdir(-1)
		if err != nil {
			t.Errorf("Readdir: unexpected error %v", err)
		}
		if len(entries) < 2 {
			t.Errorf("Readdir: expected at least 2 entries, got %d", len(entries))
		}

		names := make(map[string]bool)
		for _, entry := range entries {
			names[entry.Name()] = true
		}
		if !names["file.txt"] {
			t.Error("Readdir: missing 'file.txt'")
		}
		if !names["subdir"] {
			t.Error("Readdir: missing 'subdir'")
		}
	})

	t.Run("Readdirnames lists directory names", func(t *testing.T) {
		dir, err := rfs.Open("/testdir")
		if err != nil {
			t.Fatal(err)
		}
		defer dir.Close()

		names, err := dir.Readdirnames(-1)
		if err != nil {
			t.Errorf("Readdirnames: unexpected error %v", err)
		}
		if len(names) < 2 {
			t.Errorf("Readdirnames: expected at least 2 names, got %d", len(names))
		}

		nameSet := make(map[string]bool)
		for _, name := range names {
			nameSet[name] = true
		}
		if !nameSet["file.txt"] {
			t.Error("Readdirnames: missing 'file.txt'")
		}
		if !nameSet["subdir"] {
			t.Error("Readdirnames: missing 'subdir'")
		}
	})

	t.Run("Close closes file without error", func(t *testing.T) {
		file, err := rfs.Open("/testdir/file.txt")
		if err != nil {
			t.Fatal(err)
		}

		err = file.Close()
		if err != nil {
			t.Errorf("Close: unexpected error %v", err)
		}
	})

	t.Run("Sync succeeds (no-op for read-only)", func(t *testing.T) {
		file, err := rfs.Open("/testdir/file.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		err = file.Sync()
		if err != nil {
			t.Errorf("Sync: unexpected error %v", err)
		}
	})
}

// Phase 3: Test FileSystem Read Operations

func TestFileSystemReadOperations(t *testing.T) {
	rfs, _ := setupTestFS(t)

	t.Run("Open opens file successfully", func(t *testing.T) {
		file, err := rfs.Open("/testdir/file.txt")
		if err != nil {
			t.Errorf("Open: unexpected error %v", err)
		}
		if file == nil {
			t.Error("Open: expected non-nil file")
		} else {
			file.Close()
		}
	})

	t.Run("Stat returns FileInfo for existing files", func(t *testing.T) {
		info, err := rfs.Stat("/testdir/file.txt")
		if err != nil {
			t.Errorf("Stat: unexpected error %v", err)
		}
		if info.Name() != "file.txt" {
			t.Errorf("Stat: expected name 'file.txt', got %q", info.Name())
		}
	})

	t.Run("Lstat returns FileInfo for symlinks", func(t *testing.T) {
		info, err := rfs.Lstat("/testdir/link.txt")
		if err != nil {
			t.Errorf("Lstat: unexpected error %v", err)
		}
		if info == nil {
			t.Error("Lstat: expected non-nil FileInfo")
		}
	})

	t.Run("Readlink reads symlink target", func(t *testing.T) {
		target, err := rfs.Readlink("/testdir/link.txt")
		if err != nil {
			t.Errorf("Readlink: unexpected error %v", err)
		}
		if target != "/testdir/file.txt" {
			t.Errorf("Readlink: expected '/testdir/file.txt', got %q", target)
		}
	})

	t.Run("TempDir returns temp directory", func(t *testing.T) {
		tmpDir := rfs.TempDir()
		if tmpDir == "" {
			t.Error("TempDir: expected non-empty string")
		}
	})

	t.Run("Getwd returns current working directory", func(t *testing.T) {
		wd, err := rfs.Getwd()
		if err != nil {
			t.Errorf("Getwd: unexpected error %v", err)
		}
		if wd == "" {
			t.Error("Getwd: expected non-empty string")
		}
	})

	t.Run("Chdir changes working directory", func(t *testing.T) {
		err := rfs.Chdir("/testdir")
		if err != nil {
			t.Errorf("Chdir: unexpected error %v", err)
		}

		wd, err := rfs.Getwd()
		if err != nil {
			t.Errorf("Getwd after Chdir: unexpected error %v", err)
		}
		if wd != "/testdir" {
			t.Errorf("Getwd after Chdir: expected '/testdir', got %q", wd)
		}

		// Change back
		rfs.Chdir("/")
	})
}

// Phase 3: Test Read Operations with Various File Types

func TestReadOperationsWithVariousFileTypes(t *testing.T) {
	rfs, _ := setupTestFS(t)

	t.Run("Reading regular files", func(t *testing.T) {
		data, err := ioutil.ReadFile(rfs, "/testdir/file.txt")
		if err != nil {
			t.Errorf("ReadFile: unexpected error %v", err)
		}
		if !bytes.Equal(data, []byte("test content")) {
			t.Errorf("ReadFile: expected 'test content', got %q", string(data))
		}
	})

	t.Run("Reading directories", func(t *testing.T) {
		dir, err := rfs.Open("/testdir")
		if err != nil {
			t.Errorf("Open dir: unexpected error %v", err)
		}
		defer dir.Close()

		info, err := dir.Stat()
		if err != nil {
			t.Errorf("Stat dir: unexpected error %v", err)
		}
		if !info.IsDir() {
			t.Error("expected directory")
		}
	})

	t.Run("Reading symlinks", func(t *testing.T) {
		target, err := rfs.Readlink("/testdir/link.txt")
		if err != nil {
			t.Errorf("Readlink: unexpected error %v", err)
		}
		if target != "/testdir/file.txt" {
			t.Errorf("Readlink: expected '/testdir/file.txt', got %q", target)
		}
	})

	t.Run("Reading empty files", func(t *testing.T) {
		data, err := ioutil.ReadFile(rfs, "/empty.txt")
		if err != nil {
			t.Errorf("ReadFile empty: unexpected error %v", err)
		}
		if len(data) != 0 {
			t.Errorf("ReadFile empty: expected 0 bytes, got %d", len(data))
		}
	})

	t.Run("Reading nested files", func(t *testing.T) {
		data, err := ioutil.ReadFile(rfs, "/testdir/subdir/nested.txt")
		if err != nil {
			t.Errorf("ReadFile nested: unexpected error %v", err)
		}
		if !bytes.Equal(data, []byte("nested content")) {
			t.Errorf("ReadFile nested: expected 'nested content', got %q", string(data))
		}
	})
}

// Phase 3: Test Error Handling for Read Operations

func TestErrorHandlingForReadOperations(t *testing.T) {
	rfs, _ := setupTestFS(t)

	t.Run("Opening non-existent file returns error", func(t *testing.T) {
		_, err := rfs.Open("/nonexistent.txt")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("Stat on non-existent file returns error", func(t *testing.T) {
		_, err := rfs.Stat("/nonexistent.txt")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("Lstat on non-existent file returns error", func(t *testing.T) {
		_, err := rfs.Lstat("/nonexistent.txt")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("Readlink on non-existent path returns error", func(t *testing.T) {
		_, err := rfs.Readlink("/nonexistent_link")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("Chdir to non-existent directory returns error", func(t *testing.T) {
		err := rfs.Chdir("/nonexistent")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

// Phase 4: Edge Case Tests

func TestEdgeCases(t *testing.T) {
	rfs, _ := setupTestFS(t)

	t.Run("Operations on root directory", func(t *testing.T) {
		dir, err := rfs.Open("/")
		if err != nil {
			t.Errorf("Open root: unexpected error %v", err)
		}
		defer dir.Close()

		entries, err := dir.Readdir(-1)
		if err != nil {
			t.Errorf("Readdir root: unexpected error %v", err)
		}
		if len(entries) == 0 {
			t.Error("expected non-empty root directory")
		}
	})

	t.Run("Seek with different whence values", func(t *testing.T) {
		file, err := rfs.Open("/testdir/file.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		// SeekStart
		pos, err := file.Seek(5, io.SeekStart)
		if err != nil || pos != 5 {
			t.Errorf("SeekStart: expected 5, got %d, err %v", pos, err)
		}

		// SeekCurrent
		pos, err = file.Seek(2, io.SeekCurrent)
		if err != nil || pos != 7 {
			t.Errorf("SeekCurrent: expected 7, got %d, err %v", pos, err)
		}

		// SeekEnd
		pos, err = file.Seek(-4, io.SeekEnd)
		if err != nil {
			t.Errorf("SeekEnd: unexpected error %v", err)
		}
		// File is "test content" (12 bytes), -4 from end = 8
		if pos != 8 {
			t.Errorf("SeekEnd: expected 8, got %d", pos)
		}
	})

	t.Run("ReadAt with zero offset", func(t *testing.T) {
		file, err := rfs.Open("/testdir/file.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		buf := make([]byte, 4)
		n, err := file.ReadAt(buf, 0)
		if err != nil && err != io.EOF {
			t.Errorf("ReadAt offset 0: unexpected error %v", err)
		}
		if !bytes.Equal(buf[:n], []byte("test")) {
			t.Errorf("ReadAt offset 0: expected 'test', got %q", string(buf[:n]))
		}
	})

	t.Run("Multiple sequential reads", func(t *testing.T) {
		file, err := rfs.Open("/testdir/file.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		buf1 := make([]byte, 4)
		n1, _ := file.Read(buf1)

		buf2 := make([]byte, 4)
		n2, _ := file.Read(buf2)

		if !bytes.Equal(buf1[:n1], []byte("test")) {
			t.Errorf("First read: expected 'test', got %q", string(buf1[:n1]))
		}
		if !bytes.Equal(buf2[:n2], []byte(" con")) {
			t.Errorf("Second read: expected ' con', got %q", string(buf2[:n2]))
		}
	})
}

// Phase 4: Integration Tests

func TestCompleteWorkflows(t *testing.T) {
	rfs, _ := setupTestFS(t)

	t.Run("Complete file workflow: open -> read -> seek -> read -> close", func(t *testing.T) {
		file, err := rfs.Open("/testdir/file.txt")
		if err != nil {
			t.Fatal(err)
		}

		// First read
		buf := make([]byte, 4)
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			t.Errorf("First read: unexpected error %v", err)
		}
		if !bytes.Equal(buf[:n], []byte("test")) {
			t.Errorf("First read: expected 'test', got %q", string(buf[:n]))
		}

		// Seek back to start
		pos, err := file.Seek(0, io.SeekStart)
		if err != nil || pos != 0 {
			t.Errorf("Seek: expected 0, got %d, err %v", pos, err)
		}

		// Read all
		allBuf := make([]byte, 100)
		n, err = file.Read(allBuf)
		if err != nil && err != io.EOF {
			t.Errorf("Second read: unexpected error %v", err)
		}
		if !bytes.Equal(allBuf[:n], []byte("test content")) {
			t.Errorf("Second read: expected 'test content', got %q", string(allBuf[:n]))
		}

		// Close
		err = file.Close()
		if err != nil {
			t.Errorf("Close: unexpected error %v", err)
		}
	})

	t.Run("Directory traversal operations", func(t *testing.T) {
		// List root
		root, err := rfs.Open("/")
		if err != nil {
			t.Fatal(err)
		}
		rootEntries, err := root.Readdir(-1)
		root.Close()
		if err != nil {
			t.Errorf("Readdir root: unexpected error %v", err)
		}

		// Find testdir
		foundTestdir := false
		for _, entry := range rootEntries {
			if entry.Name() == "testdir" && entry.IsDir() {
				foundTestdir = true
				break
			}
		}
		if !foundTestdir {
			t.Error("testdir not found in root")
		}

		// List testdir
		testdir, err := rfs.Open("/testdir")
		if err != nil {
			t.Fatal(err)
		}
		testdirEntries, err := testdir.Readdir(-1)
		testdir.Close()
		if err != nil {
			t.Errorf("Readdir testdir: unexpected error %v", err)
		}

		// Check expected entries
		names := make(map[string]bool)
		for _, entry := range testdirEntries {
			names[entry.Name()] = true
		}
		if !names["file.txt"] {
			t.Error("file.txt not found in testdir")
		}
		if !names["subdir"] {
			t.Error("subdir not found in testdir")
		}
	})

	t.Run("Nested directory structures", func(t *testing.T) {
		// Read nested file
		data, err := ioutil.ReadFile(rfs, "/testdir/subdir/nested.txt")
		if err != nil {
			t.Errorf("ReadFile nested: unexpected error %v", err)
		}
		if !bytes.Equal(data, []byte("nested content")) {
			t.Errorf("ReadFile nested: expected 'nested content', got %q", string(data))
		}

		// Stat nested directory
		info, err := rfs.Stat("/testdir/subdir")
		if err != nil {
			t.Errorf("Stat subdir: unexpected error %v", err)
		}
		if !info.IsDir() {
			t.Error("expected subdir to be directory")
		}
	})
}

// Test with osfs underlying filesystem
func TestWithOSFS(t *testing.T) {
	osFS, err := osfs.NewFS()
	if err != nil {
		t.Skip("osfs not available:", err)
	}

	roFS, err := rofs.NewFS(osFS)
	if err != nil {
		t.Fatal(err)
	}

	// Verify we can't create files
	_, err = roFS.Create("/tmp/test_rofs_creation.txt")
	if !errors.Is(err, os.ErrPermission) {
		t.Errorf("Create with osfs: expected os.ErrPermission, got %v", err)
	}

	// Verify we can read existing directory
	entries, err := roFS.Open("/")
	if err != nil {
		t.Errorf("Open root with osfs: unexpected error %v", err)
	} else {
		entries.Close()
	}
}
