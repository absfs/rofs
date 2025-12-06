package rofs_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/absfs/ioutil"
	"github.com/absfs/memfs"
	"github.com/absfs/rofs"
)

// TestWrapperSuite tests rofs using fstesting.WrapperSuite.
// Note: WrapperSuite has limitations with ReadOnly wrappers (it tries to
// create test directories even when ReadOnly is true), so we run our own
// custom wrapper tests instead.
func TestWrapperSuite(t *testing.T) {
	baseFS, err := memfs.NewFS()
	if err != nil {
		t.Fatalf("failed to create base filesystem: %v", err)
	}

	// Create test directory and files in base since rofs is read-only
	testDir := filepath.Join(baseFS.TempDir(), "rofs_wrapper_test")
	if err := baseFS.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("failed to create test directory in base: %v", err)
	}
	defer baseFS.RemoveAll(testDir)

	// Create test files in base
	testFile := filepath.Join(testDir, "test_read.txt")
	testContent := []byte("test content for read operations")
	if err := ioutil.WriteFile(baseFS, testFile, testContent, 0644); err != nil {
		t.Fatalf("failed to create test file in base: %v", err)
	}

	// Create wrapped read-only filesystem
	// memfs.FileSystem already implements SymlinkFileSystem
	roFS, err := rofs.NewFS(baseFS)
	if err != nil {
		t.Fatalf("failed to create rofs: %v", err)
	}

	t.Run("ReadOperations", func(t *testing.T) {
		// Test Open and Read
		f, err := roFS.Open(testFile)
		if err != nil {
			t.Fatalf("Open failed: %v", err)
		}
		defer f.Close()

		got, err := io.ReadAll(f)
		if err != nil {
			t.Fatalf("ReadAll failed: %v", err)
		}

		if !bytes.Equal(got, testContent) {
			t.Errorf("content mismatch: got %q, want %q", got, testContent)
		}

		// Test Stat
		info, err := roFS.Stat(testFile)
		if err != nil {
			t.Fatalf("Stat failed: %v", err)
		}
		if info.Size() != int64(len(testContent)) {
			t.Errorf("file size: got %d, want %d", info.Size(), len(testContent))
		}
	})

	t.Run("WriteBlocking", func(t *testing.T) {
		// Test that write operations are blocked
		tests := []struct {
			name string
			fn   func() error
		}{
			{"Create", func() error {
				_, err := roFS.Create(filepath.Join(testDir, "newfile.txt"))
				return err
			}},
			{"Mkdir", func() error {
				return roFS.Mkdir(filepath.Join(testDir, "newdir"), 0755)
			}},
			{"MkdirAll", func() error {
				return roFS.MkdirAll(filepath.Join(testDir, "new/nested/dir"), 0755)
			}},
			{"Remove", func() error {
				return roFS.Remove(testFile)
			}},
			{"RemoveAll", func() error {
				return roFS.RemoveAll(testDir)
			}},
			{"Rename", func() error {
				return roFS.Rename(testFile, filepath.Join(testDir, "renamed.txt"))
			}},
			{"Truncate", func() error {
				return roFS.Truncate(testFile, 0)
			}},
			{"Chmod", func() error {
				return roFS.Chmod(testFile, 0600)
			}},
			{"Chtimes", func() error {
				return roFS.Chtimes(testFile, time.Now(), time.Now())
			}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.fn()
				if err == nil {
					t.Errorf("%s should fail for read-only wrapper", tt.name)
				}
				if !os.IsPermission(err) {
					t.Errorf("%s: expected permission error, got %v", tt.name, err)
				}
			})
		}
	})

	t.Run("FileWriteBlocking", func(t *testing.T) {
		f, err := roFS.Open(testFile)
		if err != nil {
			t.Fatalf("Open failed: %v", err)
		}
		defer f.Close()

		// Test Write operations on file
		_, err = f.Write([]byte("should fail"))
		if err == nil {
			t.Error("Write should fail for read-only file")
		}

		_, err = f.WriteString("should fail")
		if err == nil {
			t.Error("WriteString should fail for read-only file")
		}

		err = f.Truncate(0)
		if err == nil {
			t.Error("Truncate should fail for read-only file")
		}
	})
}
