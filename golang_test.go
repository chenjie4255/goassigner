package main

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"
)

func TestRegExpMatch(t *testing.T) {
	reg := regexp.MustCompile(parseReg)
	text := "@goassigner:As:github.com/chenjie4255/12"
	result := reg.FindStringSubmatch(text)
	if len(result) != 3 {
		t.Fatal("should match")
	}

	if result[1] != "As" {
		t.Fatalf("capture fail, wanted As, got:%s", result[1])
	}

	if result[2] != "github.com/chenjie4255/12" {
		t.Fatalf("capture fail, wanted github.com/chenjie4255/12, got:%s", result[2])
	}

	text = "@goassigner:As"
	result = reg.FindStringSubmatch(text)
	fmt.Println(result)
	if len(result) != 3 {
		t.Fatalf("should match, wanted:3 got:%d", len(result))
	}

	if result[1] != "As" {
		t.Fatalf("capture fail, wanted As, got:%s", result[1])
	}
}

func TestStructTagParse(t *testing.T) {
	ss := `gas:"123"`
	tag := reflect.StructTag(ss)
	gas := tag.Get("gas")
	if gas != "123" {
		t.Fatal("not cool")
	}
}
