package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"unicode"
)

type purser interface {
	Get(string) (string, bool)
}

type purse struct {
	items map[string][]item
}

type item struct {
	Name  string
	Query string
}

func newPurse(dir string) (p *purse, err error) {
	p = &purse{items: make(map[string][]item)}
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
			p.items[path2Key(path)] = append(p.items[path2Key(path)],
				item{Name: path2Name(path),
					Query: "`" + string(b) + "`",
				},
			)
		}
		return nil
	})

	return p, err
}

func (p *purse) getContents(filename string) (v []item, ok bool) {
	v, ok = p.items[filename]
	return
}

func (p *purse) files() []string {
	fs := make([]string, len(p.items))
	i := 0
	for k := range p.items {
		fs[i] = k
		i++
	}
	return fs
}

func path2Key(p string) string {
	return camel(filepath.Dir(p), true)
}

func path2Name(p string) string {
	return camel(filepath.Base(p), true)
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
	}

	return string(buf)
}
