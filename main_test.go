package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
)

const (
	testToteDir = "test_tote/sqltote"
	testRefDir  = "test_ref"
)

func TestOptions(t *testing.T) {
	opts := newOptions()
	if err := opts.validate(); err == nil {
		t.Fatal("error not encountered when lacking necessary info")
	}
	if err := os.Setenv("GOPACKAGE", "main"); err != nil {
		t.Fatal(err)
	}
	if err := opts.validate(); err != nil {
		t.Fatal("unexpected error during validation (setenv)")
	}

	opts.out = testToteDir
	if err := opts.validate(); err != nil {
		t.Fatal("unexpected error during validation (custom out)")
	}

	opts.pkg = testRefDir
	if err := opts.validate(); err != nil {
		t.Fatal("unexpected error during validation (pkg set)")
	}

	if err := os.Unsetenv("GOPACKAGE"); err != nil {
		t.Fatal(`Could not unset "GOPACKAGE" envar`)
	}
}

func TestMultiple(t *testing.T) {
	opts0 := newOptions()
	opts0.pkg = "main"
	if err := opts0.validate(); err != nil {
		t.Fatal("unexpected error during validation (multi simple)")
	}
	if err := mainSub(opts0); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := cleanup(opts0.defFile); err != nil {
			t.Fatal(err)
		}
	}()

	opts1 := newOptions()
	opts1.out = testToteDir
	if err := opts1.validate(); err != nil {
		t.Fatal("unexpected error during validation (multi complex no prefix)")
	}
	if err := mainSub(opts1); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := cleanup(filepath.Dir(testToteDir)); err != nil {
			t.Fatal(err)
		}
	}()

	opts2 := newOptions()
	opts2.out = testToteDir
	opts2.file = "sqltote_extra.go"
	opts2.prefix = "extra"
	if err := opts2.validate(); err != nil {
		t.Fatal("unexpected error during validation (multi complex w/prefix)")
	}
	if err := mainSub(opts2); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		fRef string
		fOut string
	}{
		{testRefDir + "/sqltote_root.go.test", testToteDir + "/sqltote.go"},
		{testRefDir + "/sqltote_extra.go.test", testToteDir + "/sqltote_extra.go"},
		{testRefDir + "/sqltote_main.go.test", "sqltote.go"},
	}

	for _, v := range tests {
		ok, err := compareFiles(v.fRef, v.fOut)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Errorf(`Generated file "%v" not equal to reference "%v".`,
				v.fOut, v.fRef,
			)
		}
	}
}

func TestFSExtended(t *testing.T) {
	opts := newOptions()
	opts.pkg = "main"
	opts.in = "sqltote/empty"
	if err := opts.validate(); err != nil {
		t.Fatal("unexpected error during validation (FS setup 1)")
	}
	if err := mainSub(opts); err == nil {
		t.Fatal("Expected error (empty file is read as dir)")
	}

	defer func() {
		if err := cleanup(opts.defFile); err != nil {
			t.Fatal(err)
		}
	}()

	opts.in = ""
	if err := opts.validate(); err != nil {
		t.Fatal("unexpected error during validation (FS setup 2)")
	}
	if err := mainSub(opts); err == nil {
		t.Fatal("Expected error (no file is accepted)")
	}

	opts.out = "sqltote/empty"
	opts.in = opts.defIn
	if err := opts.validate(); err != nil {
		t.Fatal("unexpected error during validation (FS setup 3)")
	}
	if err := mainSub(opts); err == nil {
		t.Fatal("Expected error (directory overwrites file)")
	}

	opts.out = "sqltote"
	opts.file = "bad_dir"
	if err := opts.validate(); err != nil {
		t.Fatal("unexpected error during validation (FS setup 4)")
	}
	if err := mainSub(opts); err == nil {
		t.Fatal("Expected error (file overwrites directory)")
	}
}

func compareFiles(filepath1, filepath2 string) (bool, error) {
	chunkSize := 4096
	f1, err := os.Open(filepath1)
	if err != nil {
		return false, err
	}

	f2, err := os.Open(filepath2)
	if err != nil {
		return false, err
	}

	for {
		b1 := make([]byte, chunkSize)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, chunkSize)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true, nil
			} else if err1 == io.EOF || err2 == io.EOF {
				return false, nil
			} else {
				return false, errors.New(err1.Error() + " " + err2.Error())
			}
		}

		if !bytes.Equal(b1, b2) {
			return false, nil
		}
	}
}

func cleanup(path string) error {
	if path == "./" || path == "" {
		return errors.New("Cannot remove project directory.")
	}
	if _, err := os.Stat(path); err == nil {
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}
	return nil
}
