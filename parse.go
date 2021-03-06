package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"regexp"
	"strings"
)

var importPath map[string]bool

func init() {
	importPath = make(map[string]bool)
}

type assignObject struct {
	TopStructType string
	FieldPrefix   []string
	Fields        []string
	LinkPackage   string
	LinkObject    string
}

func parseStructFields(topStructType string, fieldPrefix []string, st *ast.StructType) ([]string, []assignObject) {
	objectFields := []assignObject{}
	normalFields := []string{}
	for _, f := range st.Fields.List {
		// check if field is a embeded struct
		_, isStruct := f.Type.(*ast.StructType)
		if isStruct {
			fmt.Println("found a embed struct ")
		}
		if isStruct && len(f.Names) > 0 { // len > 0 indicate this is not an anymouse field
			// append field(struct) 's name
			newFieldPrefix := append(fieldPrefix, f.Names[0].Name)
			objs := parseEmbedStruct(topStructType, newFieldPrefix, f)
			if len(objs) > 0 {
				objectFields = append(objectFields, objs...)
			}
			continue
		}

		// if field is not a struct, just record field's name
		if f.Tag == nil {
			continue
		}

		trueTag := f.Tag.Value[1 : len(f.Tag.Value)-1]
		tag := reflect.StructTag(trueTag)
		gas := tag.Get("gas")
		if gas != "" {
			for _, n := range f.Names {
				normalFields = append(normalFields, n.Name)
			}
		}
	}

	return normalFields, objectFields
}

func parseEmbedStruct(topStructType string, fieldPrefix []string, field *ast.Field) (ret []assignObject) {
	fmt.Println("parse embed struct")
	object := &assignObject{}
	object.FieldPrefix = fieldPrefix
	object.TopStructType = topStructType

	if !parseComment(object, field.Doc) {
		return nil
	}

	structType, ok := field.Type.(*ast.StructType)
	if !ok {
		panic("wth??")
	}

	fields, embedObjs := parseStructFields(topStructType, fieldPrefix, structType)
	object.Fields = fields
	ret = append(ret, embedObjs...)
	ret = append(ret, *object)

	return
}

func parseComment(object *assignObject, commentGroup *ast.CommentGroup) bool {
	if commentGroup == nil {
		return false
	}

	for _, comment := range commentGroup.List {
		linkObjectName, packagePath := parseAssignerComment(comment.Text)
		if linkObjectName != "" {
			object.LinkObject = linkObjectName
			if packagePath != "" {
				paths := strings.Split(packagePath, "/")
				object.LinkPackage = paths[len(paths)-1]
				importPath[packagePath] = true
			}
			return true
		} else {
			fmt.Printf("parse comment(%s) fail\n", comment.Text)
		}
	}

	return false

}

func parseGeneDecl(genDecl *ast.GenDecl) (ret []assignObject) {
	// check comment
	object := &assignObject{}
	fmt.Printf("parse for %vl\n", genDecl)

	if !parseComment(object, genDecl.Doc) {
		fmt.Printf("no comment found...\n")
		return nil
	}

	for _, spec := range genDecl.Specs {
		if typeSpec, ok := spec.(*ast.TypeSpec); ok {
			fmt.Printf("found struct:%s\n", typeSpec.Name.Name)
			if typeSpec.Name != nil {
				object.TopStructType = typeSpec.Name.Name
			}

			structDecl, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				break
			}

			fields, embedObjs := parseStructFields(object.TopStructType, []string{}, structDecl)
			object.Fields = fields
			if len(embedObjs) > 0 {
				ret = append(ret, embedObjs...)
			}
		}
	}

	if object.TopStructType != "" && object.LinkObject != "" {
		ret = append(ret, *object)
	}

	return
}

func parseFile(path string) (string, []assignObject, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return "", nil, err
	}

	packageName, err := parsePackageName(f)
	if err != nil {
		return "", nil, err
	}

	ret := []assignObject{}

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			fmt.Printf("%v not a GenDecl.\n", decl)
			continue
		}

		objs := parseGeneDecl(genDecl)
		if len(objs) > 0 {
			ret = append(ret, objs...)
		}
	}

	return packageName, ret, nil
}

func parsePackageName(f *ast.File) (string, error) {
	if f.Name == nil {
		return "", errors.New("no package name found")
	}
	return f.Name.Name, nil
}

const (
	parseReg = `@goassigner:([A-Z][a-z0-9A-Z]+)[:]{0,1}([A-Za-z0-9/\-_\.]*)`
)

func parseAssignerComment(text string) (linkName, linkPackagePath string) {
	reg := regexp.MustCompile(parseReg)
	result := reg.FindStringSubmatch(text)
	if len(result) == 3 {
		linkName = result[1]
		linkPackagePath = result[2]
	} else {
		fmt.Printf("parse reg text(%s) fail, result:%v\n", text, result)
	}

	return
}
