package main

import (
	"github.com/emicklei/proto"
	"testing"
)

func TestIsScalarType(t *testing.T) {
	for _, scalarType := range scalarTypes {
		if !isScalarType(scalarType) {
			t.Errorf("isScalarType(%s) = false, want true", scalarType)
		}
	}

	if isScalarType("MyMessage") {
		t.Errorf("isScalarType(%s) = true, want false", "not a scalar type")
	}
}

func TestFullQualifiedName(t *testing.T) {
	tests := []struct {
		pkgName string
		name    string
		want    string
	}{
		{"", "MyMessage", "MyMessage"},
		{"my.package", "MyMessage", "my.package.MyMessage"},
		{"my.package", "uint32", "uint32"},
	}

	for _, test := range tests {
		got := fullQualifiedName(&test.pkgName, test.name)
		if got != test.want {
			t.Errorf("fullQualifiedName(%s, %s) = %s, want %s", test.pkgName, test.name, got, test.want)
		}
	}
}

func TestComment(t *testing.T) {
	tests := []struct {
		comment *proto.Comment
		want    string
	}{
		{nil, ""},
		{&proto.Comment{Lines: []string{"line1", "", "line2"}}, "<p>line1</p>\n"},
	}

	for _, test := range tests {
		got := string(comment(test.comment))
		if got != test.want {
			t.Errorf("comment(%v) = %s, want %s", test.comment, got, test.want)
		}
	}
}

func TestFullComment(t *testing.T) {
	tests := []struct {
		comment *proto.Comment
		want    string
	}{
		{nil, ""},
		{&proto.Comment{Lines: []string{"line1", "", "line2"}}, "<p>line1</p>\n<p>line2</p>\n"},
	}

	for _, test := range tests {
		got := string(fullComment(test.comment))
		if got != test.want {
			t.Errorf("fullComment(%v) = %s, want %s", test.comment, got, test.want)
		}
	}
}
