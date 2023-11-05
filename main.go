package main

import (
	"flag"
	"github.com/emicklei/proto"
	"log"
	"net/http"
	"os"
)

// runServer runs a http server serving the generated documentation on the given address.
func runServer(addr string, doc *string) error {
	log.Printf("Starting http server on %s", addr)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(*doc))
		if err != nil {
			log.Print(err)
		}
	})

	return http.ListenAndServe(addr, nil)
}

func main() {
	outputFlag := flag.String("out", "", "output file (default: stdout)")
	httpFlag := flag.String("http", "", "run http server (default: off) - example: -http=:8000")
	// TODO styleFlag := flag.String("style", "", "custom css style file")
	flag.Parse()

	var definitions []*proto.Proto
	if len(flag.Args()) == 0 {
		// Read from stdin if no files are given.
		parser := proto.NewParser(os.Stdin)
		definition, err := parser.Parse()
		if err != nil {
			log.Fatal(err)
		}

		definitions = append(definitions, definition)
	} else {
		for _, arg := range flag.Args() {
			file, err := os.Open(arg)
			if err != nil {
				log.Fatal(err)
			}

			parser := proto.NewParser(file)
			definition, err := parser.Parse()
			if err != nil {
				log.Fatal(err)
			}

			definitions = append(definitions, definition)

			err = file.Close()
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	doc, err := GenerateDoc(definitions...)
	if err != nil {
		log.Fatal(err)
	}

	if *httpFlag != "" {
		log.Fatal(runServer(*httpFlag, &doc))

		return
	}

	if *outputFlag == "" {
		_, err := os.Stdout.WriteString(doc)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		file, err := os.Create(*outputFlag)
		if err != nil {
			log.Fatal(err)
		}

		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(file)

		_, err = file.WriteString(doc)
		if err != nil {
			log.Fatal(err)
		}
	}
}
