package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type purse interface {
	Get(string) (string, bool)
}

type memoryPurse struct {
	mu sync.RWMutex
	fs map[string]string
}

func newPurse(dir string) (*memoryPurse, error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fis, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}

	p := &memoryPurse{fs: make(map[string]string, 0)}

	for _, fi := range fis {
		if !fi.IsDir() && filepath.Ext(fi.Name()) == ext {
			f, err := os.Open(filepath.Join(dir, fi.Name()))
			if err != nil {
				return nil, err
			}
			b, err := ioutil.ReadAll(f)
			if err != nil {
				return nil, err
			}
			p.fs[fi.Name()] = string(b)
			f.Close()
		}
	}
	return p, nil
}

func (p *memoryPurse) getContents(filename string) (v string, ok bool) {
	p.mu.RLock()
	v, ok = p.fs[filename]
	p.mu.RUnlock()
	return
}

func (p *memoryPurse) files() []string {
	fs := make([]string, len(p.fs))
	i := 0
	for k, _ := range p.fs {
		fs[i] = k
		i++
	}
	return fs
}
