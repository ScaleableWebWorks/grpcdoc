package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/emicklei/proto"
	"html/template"
	"strings"
)

//go:embed internal/doc.html
var doc string

//go:embed internal/style.css
var style string

var scalarTypes = [...]string{"double", "float", "int32", "int64", "uint32", "uint64",
	"sint32", "sint64", "fixed32", "fixed64", "sfixed32", "sfixed64", "bool", "string", "bytes"}

type ServiceDoc struct {
	Name        string
	Description string
	Methods     []*MethodDoc
}

type MessageDoc struct {
	Name        string
	Description string
	Fields      []*FieldDoc
}

type FieldDoc struct {
	Position    uint8
	Name        string
	Description string
	Type        string
	IsScalar    bool
	IsRepeated  bool
}

type MethodDoc struct {
	Name        string
	Description string
	Input       string
	Output      string
}

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

func GenerateDoc(definitions ...*proto.Proto) (string, error) {
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
							Position:    uint8(field.Sequence),
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
							Position:    uint8(field.Integer),
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

	data := Data{
		Filename: "",
		Style:    template.CSS(style),
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

func isScalarType(t string) bool {
	for _, scalarType := range scalarTypes {
		if t == scalarType {
			return true
		}
	}
	return false
}

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

func comment(c *proto.Comment) string {
	if c == nil {
		return ""
	}

	return c.Message()
}
