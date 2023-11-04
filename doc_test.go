package main

import (
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
