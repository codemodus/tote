// Tote is a CLI application for generating structs which store SQL queries
// as defined by the directory and file structure supplied (default is
// "./sqltote").  Only .sql files are read.
//
// 	Available flags:
// 	--in={dir}          Set the SQL storage directory.
// 	--out={dir}         Set the tote package directory.
// 	--file={filename}   Set the tote file name.
// 	--pkg={package}     Set the tote package name.
// 	--prefix={name}     Set the tote struct prefix.
//
// Normally, this command should be called using go:generate.  The following
// usage will produce a package named "totepkg" within the "totepkg"
// directory:
// 	//go:generate tote -in=resources/sql/tote -out=totepkg
//
// The following usage will add a second file to the "totepkg" package:
// 	//go:generate tote -in=other/sql/tote -out=totepkg -prefix=other -file=other.go
//
// Queries are accessible in this way:
// 	import "vcs-storage.nil/mycurrentproject/totepkg"
//
// 	func main() {
// 		// File originally located at "./resources/sql/tote/user/all.sql"
// 		fmt.Println(totepkg.User.All)
//
// 		// File originally located at "./resources/sql/tote/user/role/many_by_user.sql"
// 		fmt.Println(totepkg.UserRole.ManyByUser)
//
// 		// File originally located at "./other/sql/tote/user/one_by_name.sql"
// 		fmt.Println(totepkg.OtherUser.OneByName)
// 	}
//
// The main caveat seems to be naming collisions which was the primary
// motivation for the prefix flag.  Stay aware and problems can be avoided.
//
// This package started as a fork of smotes/purse.
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

func main() {
	const (
		defFile = "sqltote.go"
		defIn   = "sqltote"
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

	f, err := os.Create(filepath.Join(out, file))
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	t, err := newTote(in, prefix)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := &tmplContext{
		Pkg:   pkg,
		Items: t.Items,
	}

	b := &bytes.Buffer{}
	pt, err := template.New(defIn).Parse(tmpl)
	if err != nil {
		log.Fatalln(err)
	}
	if err = pt.Execute(b, ctx); err != nil {
		log.Fatalln(err)
	}

	fb, err := format.Source(b.Bytes())
	if err != nil {
		log.Fatalln(err)
	}

	f.Write(fb)
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
