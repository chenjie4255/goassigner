package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"
)

func main() {
	inputFile := ""
	flag.StringVar(&inputFile, "f", "", "output file")
	flag.Parse()

	packageName, objs, err := parseFile(inputFile)
	if err != nil {
		fmt.Printf("parse file error:%s", err)
		return
	}

	fmt.Printf("parse result:%+v", objs)
	outputPath := getOutputPath(inputFile)
	render(outputPath, packageName, objs)
}

func getOutputPath(inputPath string) string {
	if !strings.HasSuffix(inputPath, ".go") {
		return ""
	}
	dir, file := filepath.Split(inputPath[:len(inputPath)-3])
	return filepath.Join(dir, fmt.Sprintf("%s_assigner.go", file))
}
