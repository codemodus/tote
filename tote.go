package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type tote struct {
	Items map[string][]item
}

type item struct {
	Name  string
	Query string
}

func newTote(dir, prefix string) (t *tote, err error) {
	t = &tote{Items: make(map[string][]item)}
	if _, err = os.Stat(dir); err != nil {
		return nil, err
	}

	err = filepath.Walk(dir, func(path string, i os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		if !i.IsDir() && filepath.Ext(path) == ".sql" {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			b, err := ioutil.ReadAll(f)
			if err != nil {
				return err
			}
			t.Items[path2Key(path, dir, prefix)] = append(
				t.Items[path2Key(path, dir, prefix)],
				item{Name: path2Name(path),
					Query: "`" + string(b) + "`",
				},
			)
		}
		return nil
	})

	return t, err
}

func path2Key(p, dir, prefix string) string {
	r := camel(strings.TrimPrefix(filepath.Dir(p), dir), true)
	if prefix != "" {
		r = camel(prefix, true) + r
	}
	if r == "" {
		r = "Root"
	}
	return r
}

func path2Name(p string) string {
	return camel(strings.TrimSuffix(filepath.Base(p), ".sql"), true)
}

func camel(s string, ucFirst bool) string {
	rs := []rune(s)
	l := len(rs)
	buf := make([]rune, 0, l)

	for i := 0; i < l; i++ {
		if unicode.IsLetter(rs[i]) {
			if i == 0 && ucFirst || i > 0 && !unicode.IsLetter(rs[i-1]) {
				buf = append(buf, unicode.ToUpper(rs[i]))
			} else {
				buf = append(buf, rs[i])
			}
		}
		if unicode.IsNumber(rs[i]) {
			buf = append(buf, rs[i])
		}
	}

	return string(buf)
}
