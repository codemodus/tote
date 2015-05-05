package main

type tmplContext struct {
	Pkg string
	Fs  map[string][]item
}

const (
	tmpl = `package {{.Pkg}}

var (
{{range $key, $val := .Fs}}	// {{$key}} is something or another.
	{{$key}} = struct {
	{{range $k, $v := $val}}	{{$v.Name}} string
	{{end}}}{
	{{range $k, $v := $val}}	{{$v.Query}},
	{{end}}}
{{end}})
`
)
