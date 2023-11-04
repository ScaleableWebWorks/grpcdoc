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

type Data struct {
	Filename string
	Style    template.CSS
	Services []*ServiceDoc
	Messages []*MessageDoc
}

func GenerateDoc(definition *proto.Proto) (string, error) {
	t, err := template.New("doc").Parse(doc)
	if err != nil {
		return "", err
	}

	var services []*ServiceDoc
	var messages []*MessageDoc
	for _, element := range definition.Elements {
		switch element.(type) {
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
		}
	}

	data := Data{
		Filename: "test.proto",
		Style:    template.CSS(style),
		Services: services,
		Messages: messages,
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
