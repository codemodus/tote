package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

func main() {
	const (
		envar   = "GOPACKAGE"
		defFile = "sqlpurse.go"
		defOut  = "./"
	)
	var in, out, file, name, pkg string

	flag.StringVar(&in, "in", "", "directory of the input SQL file(s)")
	flag.StringVar(&out, "out", defOut, "directory of the output source file")
	flag.StringVar(&file, "file", defFile, "name of the output source file")
	flag.StringVar(&name, "name", "gen", "variable name of the generated Purse struct")
	flag.StringVar(&pkg, "pkg", "", "name of the go package for the generated source file")
	flag.Parse()

	if in == "" {
		log.Fatalln(errors.New("error: SQL directory must be provided"))
	}
	if pkg == "" {
		if os.Getenv(envar) == "" {
			log.Fatalln(errors.New("error: package name receiving " +
				"generated source must be provided"))
		}
		pkg = os.Getenv(envar)
	}

	p, err := newPurse(in)
	if err != nil {
		log.Fatalln(err)
	}

	fs := make(map[string][]item, len(p.files()))
	for _, name := range p.files() {
		s, ok := p.getContents(name)
		if !ok {
			log.Printf("Unable to get file %s", name)
			continue
		}
		fs[name] = s
	}

	ctx := &tmplContext{
		Pkg: pkg,
		Fs:  fs,
	}

	t, err := template.New(name).Parse(tmpl)
	if err != nil {
		log.Fatalln(err)
	}

	f, err := os.Create(filepath.Join(out, file))
	if err != nil {
		log.Fatalln(err)
	}

	err = t.Execute(f, ctx)
	if err != nil {
		log.Fatalln(err)
	}
}
