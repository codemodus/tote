package main

type tmplContext struct {
	Pkg   string
	Items map[string][]item
}

const (
	tmpl = `package {{.Pkg}}

var (
{{range $key, $val := .Items}}	// {{$key}} is a generated SQL tote.  Do not modify.
	{{$key}} = struct {
	{{range $k, $v := $val}}	{{$v.Name}} string
	{{end}}}{
	{{range $k, $v := $val}}	{{$v.Query}},
	{{end}}}
{{end}})
`
)
