package main

import (
	"flag"
	"github.com/emicklei/proto"
	"io"
	"os"
)

func main() {
	fileFlag := flag.String("file", "", "proto file to parse")
	outputFlag := flag.String("out", "", "output file (default: stdout)")

	flag.Parse()

	var reader io.Reader
	if *fileFlag == "" {
		reader = os.Stdin
	} else {
		file, err := os.Open(*fileFlag)
		if err != nil {
			panic(err)
		}

		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				panic(err)
			}
		}(file)
		reader = file
	}

	parser := proto.NewParser(reader)
	parser.Filename(*fileFlag)

	definition, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	doc, err := GenerateDoc(definition)
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
