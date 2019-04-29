package rofs_test

import (
	"bytes"
	"log"
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
		log.Fatal(err)
	}

	if bytes.Compare(dataOut, dataIn) != 0 {
		log.Fatalf("expected read to match write: %q expect %q", string(dataIn), string(dataOut))
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
	t.Logf("%s %s % 6d %s", info.Mode(), info.Name(), info.Size(), info.ModTime().Format(time.UnixDate))
}
