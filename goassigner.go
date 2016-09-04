package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

func main() {
	inputFile := ""
	flag.StringVar(&inputFile, "f", "", "output file")
	flag.Parse()

	objs, err := parseFile(inputFile)
	if err != nil {
		fmt.Printf("parse file error:%s", err)
		return
	}

	for _, o := range objs {
		fmt.Printf("ready to assign %s->%s %v\n", o.LinkObject, o.Name, o.Fields)
	}

	outputPath := getOutputPath(inputFile)
	render(outputPath, objs)
}

func getOutputPath(inputPath string) string {
	if !strings.HasSuffix(inputPath, ".go") {
		return ""
	}
	dir, file := filepath.Split(inputPath[:len(inputPath)-3])
	return filepath.Join(dir, fmt.Sprintf("%s_assigner.go", file))
}

var renderTemplate = template.Must(template.New("render").Parse(`// This file is generated by goassigner, DO NOT EDIT IT.
// (github.com/chenjie4255/goassigner) 
package {{.Package}}
{{range .Objects}}
func (s *{{.Name}}) AssignFrom{{.LinkObject}}(obj {{.LinkObject}}) {
	{{range .Fields}}s.{{.}} = obj.{{.}}
	{{end}}
}
{{end}}
`))

func render(outputPath string, objs []assignObject) {
	type RenderData struct {
		Package string
		Objects []assignObject
	}

	output, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("create/open output file(%s) fail: %s", outputPath, err)
		return
	}
	defer output.Close()

	var data RenderData
	data.Package = "example"
	data.Objects = objs

	renderTemplate.Execute(output, data)

}

func parseFile(path string) (ret []assignObject, err error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return
	}

	for _, decl := range f.Decls {
		object := &assignObject{}
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			fmt.Printf("%v not a GenDecl.\n", decl)
			continue
		}

		if genDecl.Doc == nil {
			fmt.Printf("%v doc is nil\n", decl)
			continue
		}

		found := false
		for _, comment := range genDecl.Doc.List {
			reg := regexp.MustCompile(`@goassigner:([A-Z][a-z0-9A-Z]+)`)
			result := reg.FindStringSubmatch(comment.Text)
			if len(result) != 2 {
				continue
			}

			object.LinkObject = result[1]
			found = true
			break
		}

		if !found {
			continue
		}

		for _, spec := range genDecl.Specs {
			if typeSpec, ok := spec.(*ast.TypeSpec); ok {
				if typeSpec.Name != nil {
					object.Name = typeSpec.Name.Name
				}

				structDecl, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					break
				}

				for _, f := range structDecl.Fields.List {
					for _, n := range f.Names {
						object.Fields = append(object.Fields, n.Name)
					}
				}
			}
		}

		if object.Name != "" {
			ret = append(ret, *object)
		}
	}

	return
}

type assignObject struct {
	Name       string
	Fields     []string
	LinkObject string
}
