package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/emicklei/proto"
	"html/template"
	"strings"
)

// Holds the html template for the documentation.
//
//go:embed internal/doc.html
var doc string

// Holds the default css defaultStyle for the documentation.
//
//go:embed internal/style.css
var defaultStyle string

// A list of all protobuf scalar types as defined here:
// https://developers.google.com/protocol-buffers/docs/proto3#scalar
var scalarTypes = [...]string{"double", "float", "int32", "int64", "uint32", "uint64",
	"sint32", "sint64", "fixed32", "fixed64", "sfixed32", "sfixed64", "bool", "string", "bytes"}

// ServiceDoc represents a service which is used for template rendering.
type ServiceDoc struct {
	Name        string
	Description string
	Methods     []*MethodDoc
}

// MessageDoc represents a message which is used for template rendering.
type MessageDoc struct {
	Name        string
	Description string
	Fields      []*FieldDoc
}

// FieldDoc represents a field of a message or enum type which is used for template rendering.
type FieldDoc struct {
	FieldNumber uint8
	Name        string
	Description string
	Type        string
	// IsScalar is true if the field type is a protobuf scalar type.
	IsScalar   bool
	IsRepeated bool
}

// MethodDoc represents a method of a grpc service which is used for template rendering.
type MethodDoc struct {
	Name        string
	Description string
	Input       string
	Output      string
}

// EnumDoc represents an enum which is used for template rendering.
type EnumDoc struct {
	Name        string
	Description string
	Fields      []*FieldDoc
}

type Data struct {
	Filename string
	Style    template.CSS
	Services []*ServiceDoc
	Messages []*MessageDoc
	Enums    []*EnumDoc
}

// GenerateDoc takes a list of protobuf definitions and generates a html documentation.
func GenerateDoc(customStyle *string, definitions ...*proto.Proto) (string, error) {
	t, err := template.New("doc").Parse(doc)
	if err != nil {
		return "", err
	}

	var pkgName *string = nil
	var services []*ServiceDoc
	var messages []*MessageDoc
	var enums []*EnumDoc
	for _, definition := range definitions {
		for _, element := range definition.Elements {
			switch element.(type) {
			case *proto.Import:
				imp := element.(*proto.Import)
				println("Import: ", imp.Filename)
			case *proto.Package:
				pkg := element.(*proto.Package)
				pkgName = &pkg.Name
				println("Package: ", pkg.Name)
			case *proto.Service:
				service := element.(*proto.Service)
				var methods []*MethodDoc
				for _, element := range service.Elements {
					switch element.(type) {
					case *proto.RPC:
						rpc := element.(*proto.RPC)
						doc := MethodDoc{
							Name:        rpc.Name,
							Description: comment(rpc.Comment),
							Input:       fullQualifiedName(pkgName, rpc.RequestType),
							Output:      fullQualifiedName(pkgName, rpc.ReturnsType),
						}
						methods = append(methods, &doc)
					}
				}
				doc := ServiceDoc{
					Name:        fullQualifiedName(pkgName, service.Name),
					Description: comment(service.Comment),
					Methods:     methods,
				}
				services = append(services, &doc)
			case *proto.Message:
				message := element.(*proto.Message)
				var fields []*FieldDoc
				for _, element := range message.Elements {
					switch element.(type) {
					case *proto.NormalField:
						field := element.(*proto.NormalField)

						doc := FieldDoc{
							FieldNumber: uint8(field.Sequence),
							Name:        field.Name,
							Description: comment(field.Comment),
							Type:        fullQualifiedName(pkgName, field.Type),
							IsScalar:    isScalarType(field.Type),
							IsRepeated:  field.Repeated,
						}
						fields = append(fields, &doc)
					}
				}

				doc := MessageDoc{
					Name:        fullQualifiedName(pkgName, message.Name),
					Description: comment(message.Comment),
					Fields:      fields,
				}

				messages = append(messages, &doc)
			case *proto.Enum:
				enum := element.(*proto.Enum)

				var fields []*FieldDoc
				for _, element := range enum.Elements {
					switch element.(type) {
					case *proto.EnumField:
						field := element.(*proto.EnumField)
						doc := FieldDoc{
							FieldNumber: uint8(field.Integer),
							Name:        field.Name,
							Description: comment(field.Comment),
						}
						fields = append(fields, &doc)
					}
				}
				doc := EnumDoc{
					Name:        fullQualifiedName(pkgName, enum.Name),
					Description: comment(enum.Comment),
					Fields:      fields,
				}

				enums = append(enums, &doc)
			}
		}
	}

	var style = customStyle
	if style == nil {
		style = &defaultStyle
	}

	data := Data{
		Filename: "",
		Style:    template.CSS(*style),
		Services: services,
		Messages: messages,
		Enums:    enums,
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// isScalarType returns true if the given type is a protobuf scalar type.
func isScalarType(t string) bool {
	for _, scalarType := range scalarTypes {
		if t == scalarType {
			return true
		}
	}
	return false
}

// fullQualifiedName returns the full qualified name of the given type if applicable.
// In the following cases the name is returned as is:
//   - pkgName is nil or empty
//   - name contains a dot (is already a full qualified name)
//   - name is a scalar type
func fullQualifiedName(pkgName *string, name string) string {
	if pkgName == nil || *pkgName == "" {
		return name
	} else if strings.ContainsRune(name, '.') {
		return name
	} else if isScalarType(name) {
		return name
	}

	return fmt.Sprintf("%v.%v", *pkgName, name)
}

// comment returns the comment or empty string if the comment is nil.
func comment(c *proto.Comment) string {
	if c == nil {
		return ""
	}

	return c.Message()
}
