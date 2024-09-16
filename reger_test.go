package reger

import (
	"testing"
	"strings"
	"regexp"
	"reflect"
)

const input = "here is some apple cider, and there is an apple spider"

var testRe = regexp.MustCompile(`apple (sp|c)ider`)

func TestFindReaderAllString(t *testing.T) {
	expected := []string{"apple cider", "apple spider"}
	found := NewReaderReger(strings.NewReader(input)).FindReaderAllString(testRe)
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("found %v != expected %v", found, expected)
	}
}

func TestFindReaderAllStringSubmatch(t *testing.T) {
	expected := [][]string{
		[]string{"apple cider", "c"},
		[]string{"apple spider", "sp"},
	}
	found := NewReaderReger(strings.NewReader(input)).FindReaderAllStringSubmatch(testRe)
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("found %v != expected %v", found, expected)
	}
}
