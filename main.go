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
		defFile = "sqlpurse.go"
		defIn   = "sqlpurse"
		defOut  = "./"
		envar   = "GOPACKAGE"
	)
	var file, in, out, pkg, prefix string

	flag.StringVar(&in, "in", defIn, "directory of the input SQL file(s)")
	flag.StringVar(&out, "out", defOut, "directory of the output source file")
	flag.StringVar(&file, "file", defFile, "name of the output source file")
	flag.StringVar(&pkg, "pkg", "", "name of the go package for the generated source file")
	flag.StringVar(&prefix, "prefix", "", "prefix for struct names")
	flag.Parse()

	if pkg == "" {
		if out == defOut {
			if os.Getenv(envar) == "" {
				log.Fatalln(errors.New("error: package name receiving " +
					"generated source must be provided"))
			}
			pkg = os.Getenv(envar)
		} else {
			pkg = filepath.Base(out)
		}
	}

	p, err := newPurse(in, prefix)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := &tmplContext{
		Pkg:   pkg,
		Items: p.Items,
	}

	t, err := template.New(defIn).Parse(tmpl)
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
