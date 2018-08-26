package main

import (
	"bytes"
	"errors"
	"flag"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/codemodus/kace"
)

type options struct {
	in        string
	defIn     string
	out       string
	defOut    string
	file      string
	defFile   string
	pkg       string
	defPkg    string
	prefix    string
	defPrefix string
}

type item struct {
	Name  string
	Query string
}

type tote struct {
	Items map[string][]item
}

func (opts *options) validate() error {
	const (
		envvar = "GOPACKAGE"
	)

	if opts.pkg != "" {
		return nil
	}
	if opts.out == opts.defOut {
		if os.Getenv(envvar) == "" {
			return errors.New("error: package name receiving " +
				"generated source must be provided")
		}
		opts.pkg = os.Getenv(envvar)
		return nil
	}
	opts.pkg = filepath.Base(opts.out)
	return nil
}

func path2Name(p string) string {
	return kace.Pascal(strings.TrimSuffix(filepath.Base(p), ".sql"))
}

func path2Key(p, dir, prefix string) string {
	r := kace.Pascal(strings.TrimPrefix(filepath.Dir(p), dir))
	if prefix != "" {
		r = kace.Pascal(prefix + "/" + r)
	}
	if r == "" {
		r = "Root"
	}
	return r
}

func newOptions() *options {
	const (
		defIn   = "sqltote"
		defOut  = "./"
		defFile = "sqltote.go"
	)

	return &options{
		in: defIn, defIn: defIn,
		out: defOut, defOut: defOut,
		file: defFile, defFile: defFile,
	}
}

func newTote(dir, prefix string) (t *tote, err error) {
	t = &tote{Items: make(map[string][]item)}
	fi, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, errors.New("select in directory is a file")
	}

	werr := filepath.Walk(dir, func(path string, i os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		if !i.IsDir() && filepath.Ext(path) == ".sql" {
			f, err := os.Open(path) // nolint: gosec
			if err != nil {
				return err
			}
			defer func() {
				if err = f.Close(); err != nil {
					panic(err)
				}
			}()

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
	return t, werr
}

func mainSub(opts *options) error {
	fp := filepath.Join(opts.out, opts.file)
	dir := filepath.Dir(fp)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return err
	}
	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer func() {
		if err = f.Close(); err != nil {
			panic(err)
		}
	}()

	t, err := newTote(opts.in, opts.prefix)
	if err != nil {
		return err
	}

	ctx := &tmplContext{
		Pkg:   opts.pkg,
		Items: t.Items,
	}

	b := &bytes.Buffer{}
	pt, err := template.New("sqltote").Parse(tmpl)
	if err != nil {
		panic(err)
	}
	if err = pt.Execute(b, ctx); err != nil {
		panic(err)
	}

	bs, err := format.Source(b.Bytes())
	if err != nil {
		panic(err)
	}

	if _, err := f.Write(bs); err != nil {
		panic(err)
	}
	return nil
}

func main() {
	opts := newOptions()

	flag.StringVar(&opts.in, "in", opts.defIn, "directory of the input SQL file(s)")
	flag.StringVar(&opts.out, "out", opts.defOut, "directory of the output source file")
	flag.StringVar(&opts.file, "file", opts.defFile, "name of the output source file")
	flag.StringVar(&opts.pkg, "pkg", opts.defPkg, "name of the go package for the generated source file")
	flag.StringVar(&opts.prefix, "prefix", opts.defPrefix, "prefix for struct names")
	flag.Parse()

	if err := opts.validate(); err != nil {
		log.Fatal(err)
	}

	if err := mainSub(opts); err != nil {
		log.Fatalln(err)
	}
}
