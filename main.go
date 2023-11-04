package main

import (
	"flag"
	"github.com/emicklei/proto"
	"os"
)

func main() {
	outputFlag := flag.String("out", "", "output file (default: stdout)")

	flag.Parse()

	var definitions []*proto.Proto
	for _, arg := range flag.Args() {
		file, err := os.Open(arg)
		if err != nil {
			panic(err)
		}

		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				panic(err)
			}
		}(file)

		parser := proto.NewParser(file)
		definition, err := parser.Parse()
		if err != nil {
			panic(err)
		}

		definitions = append(definitions, definition)
	}

	doc, err := GenerateDoc(definitions...)
	if err != nil {
		panic(err)
	}

	if *outputFlag == "" {
		_, err := os.Stdout.WriteString(doc)
		if err != nil {
			panic(err)
		}
	} else {
		file, err := os.Create(*outputFlag)
		if err != nil {
			panic(err)
		}

		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				panic(err)
			}
		}(file)

		_, err = file.WriteString(doc)
		if err != nil {
			panic(err)
		}
	}
}
