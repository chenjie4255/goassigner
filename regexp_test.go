package main

import (
	"regexp"
	"testing"
)

func TestRegExpMatch(t *testing.T) {
	reg := regexp.MustCompile(`@goassigner:([A-Z][a-z0-9A-Z]+)`)
	text := "@goassigner:As"
	result := reg.FindStringSubmatch(text)
	if len(result) != 2 {
		t.Fatal("should match")
	}

	if result[1] != "As" {
		t.Fatalf("capture fail, wanted As, got:%s", result[1])
	}
}
