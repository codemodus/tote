package main

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var (
	dirname  string
	fixtures = map[string]string{
		"insert.sql":        "INSERT INTO post (slug, title, created, markdown, html)\nVALUES (?, ?, ?, ?, ?)",
		"query_all.sql":     "SELECT\nid,\nslug,\ntitle,\ncreated,\nmarkdown,\nhtml\nFROM post",
		"query_by_slug.sql": "SELECT\nid,\nslug,\ntitle,\ncreated,\nmarkdown,\nhtml\nFROM post\nWHERE slug = ?",
	}
)

func init() {
	dirname = filepath.Join(".", "fixtures")

	// replace newlines for running unit tests on windows
	if runtime.GOOS == "windows" {
		for k, v := range fixtures {
			fixtures[k] = strings.Replace(v, "\n", "\r\n", -1)
		}
	}
}

func TestNew(t *testing.T) {
	s, err := newPurse(dirname)
	if err != nil {
		t.Errorf("unexpected error from New() on fixtures directory")
	}

	if len(fixtures) != len(s.items) {
		t.Errorf("invalid number of loaded SQL files")
	}

	for key := range fixtures {
		_, ok := s.items[key]
		if !ok {
			t.Errorf("unable to find loaded file %s in file map", key)
		}
	}

	// verify only SQL files were loaded
	for key := range s.items {
		if filepath.Ext(key) != ".sql" {
			t.Errorf("loaded unexpected file type: %s", key)
		}
	}

	// try to load file instead of directory
	_, err = newPurse(filepath.Join(".", "purse.go"))
	if err == nil {
		t.Errorf("expected error trying to load from non-directory")
	}

	// try to load directory that does not exist
	_, err = newPurse(filepath.Join(".", "foo"))
	if err == nil {
		t.Errorf("expected error trying to load directory that does not exist")
	}
}

func TestGet(t *testing.T) {
	s, err := newPurse(dirname)
	if err != nil {
		t.Errorf("unexpected error from New() on fixtures directory")
	}

	for key, val := range fixtures {
		v, ok := s.getContents(key)
		if !ok {
			t.Errorf("unable to find loaded file %s in file map", key)
		}
		if v != val {
			t.Errorf("invalid %s file content:\n%v\n%v", key, []byte(v), []byte(val))
		}
	}
}

func BenchmarkGet(b *testing.B) {
	s, _ := newPurse(dirname)
	key := "query_by_slug.sql"

	for i := 0; i < b.N; i++ {
		s.getContents(key)
	}
}
