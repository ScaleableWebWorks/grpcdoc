package main

import (
	"bytes"
	_ "embed"
	"github.com/emicklei/proto"
	"html/template"
)

//go:embed internal/doc.html
var doc string

//go:embed internal/style.css
var style string

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

func GenerateDoc(definition *proto.Proto) (string, error) {
	t, err := template.New("doc").Parse(doc)
	if err != nil {
		return "", err
	}

	var services []*ServiceDoc
	var messages []*MessageDoc
	var enums []*EnumDoc
	for _, element := range definition.Elements {
		switch element.(type) {
		case *proto.Import:
			imp := element.(*proto.Import)
			println("Import: ", imp.Filename)
		case *proto.Package:
			pkg := element.(*proto.Package)
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
						Description: rpc.Comment.Message(),
						Input:       rpc.RequestType,
						Output:      rpc.ReturnsType,
					}
					methods = append(methods, &doc)
				}
			}
			doc := ServiceDoc{
				Name:        service.Name,
				Description: service.Comment.Message(),
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
						Description: field.Comment.Message(),
						Type:        field.Type,
					}
					fields = append(fields, &doc)
				}
			}

			doc := MessageDoc{
				Name:        message.Name,
				Description: message.Comment.Message(),
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
						Description: field.Comment.Message(),
					}
					fields = append(fields, &doc)
				}
			}
			doc := EnumDoc{
				Name:        enum.Name,
				Description: enum.Comment.Message(),
				Fields:      fields,
			}

			enums = append(enums, &doc)
		}
	}

	data := Data{
		Filename: definition.Filename,
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
