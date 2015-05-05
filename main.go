package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

const (
	ext string = ".sql"
)

func main() {
	const envar = "GOPACKAGE"
	var in, out, file, name, pkg string

	flag.StringVar(&in, "in", "", "directory of the input SQL file(s)")
	flag.StringVar(&out, "out", "./", "directory of the output source file")
	flag.StringVar(&file, "file", "out.go", "name of the output source file")
	flag.StringVar(&name, "name", "gen", "variable name of the generated Purse struct")
	flag.StringVar(&pkg, "pkg", "", "name of the go package for the generated source file")
	flag.Parse()

	if in == "" {
		log.Fatalln(errors.New("error: SQL directory must be provided."))
	}
	if pkg == "" {
		if os.Getenv(envar) == "" {
			log.Fatalln(errors.New("error: package name receiving generated source must be provided."))
		}
		pkg = os.Getenv(envar)
	}

	p, err := newPurse(in)
	if err != nil {
		log.Fatalln(err)
	}

	fs := make(map[string]string, len(p.files()))
	for _, name := range p.files() {
		s, ok := p.getContents(name)
		if !ok {
			log.Printf("Unable to get file %s", name)
			continue
		}
		fs[name] = strconv.Quote(s)
	}

	ctx := &tmplContext{
		Varname: name,
		Package: pkg,
		Files:   fs,
	}

	cs := tmplHead + tmplBodyStruct + "\n" + tmplBodyVar

	if out != "./" {
		ctx.Varname = strings.Title(name)
		cs = tmplHead + tmplBodyVar

		tmplCommon, err := template.New(name).Parse(
			tmplHead + tmplBodyStruct)
		if err != nil {
			log.Fatalln(err)
		}

		fCmn, err := os.Create(filepath.Join(out, pkg+".go"))
		if err != nil {
			log.Fatalln(err)
		}

		err = tmplCommon.Execute(fCmn, ctx)
		if err != nil {
			log.Fatalln(err)
		}
	}

	tmpl, err := template.New(name).Parse(cs)
	if err != nil {
		log.Fatalln(err)
	}

	f, err := os.Create(filepath.Join(out, file))
	if err != nil {
		log.Fatalln(err)
	}

	err = tmpl.Execute(f, ctx)
	if err != nil {
		log.Fatalln(err)
	}
}
