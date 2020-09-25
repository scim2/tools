package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/scim2/tools/schema"
	gen "github.com/scim2/tools/structs/generate"
	"io/ioutil"
)

var (
	path    string
	out     string
	pkgName string
)

func main() {
	flag.StringVar(&path, "p", "", "path to json of the schema")
	flag.StringVar(&path, "path", "", "path to json of the schema")
	flag.StringVar(&out, "o", "", "path of generated go file")
	flag.StringVar(&out, "out", "", "path of generated go file")
	flag.StringVar(&pkgName, "pkg", "main", "package name")
	flag.Parse()

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var s schema.ReferenceSchema
	if err := json.Unmarshal(raw, &s); err != nil {
		panic(err)
	}

	g, err := gen.NewStructGenerator(s)
	if err != nil {
		panic(err)
	}

	b := &bytes.Buffer{}
	b.WriteString("// Do not edit. This file is auto-generated.\n")
	b.WriteString(fmt.Sprintf("package %s\n\n", pkgName))
	_, _ = g.Generate().WriteTo(b)

	if err := ioutil.WriteFile(out, b.Bytes(), 0644); err != nil {
		panic(err)
	}
}
