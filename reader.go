package main

import (
	"fmt"
	"os"

	"github.com/deltamobile/goraptor"
)

// type Parser struct {
// 	rdfparser *goraptor.Parser
// 	input     io.Reader
// 	index     map[string]*builder
// 	buffer    map[string][]bufferEntry
// 	doc       *spdx.Document
// }

func main() {

	// check that we've received the right number of arguments

	args := os.Args
	if len(args) != 2 {
		fmt.Printf("Usage: %v <spdx-file-in>\n", args[0])
		fmt.Printf("  Load SPDX 2.1 RDF file <spdx-file-in>, and\n")
		fmt.Printf("  print its contents.\n")
		return
	}

	//   PARSE FILE method - Takes the file location as an input
	filename := args[1]
	defer fmt.Println("DONE!")
	parserFile := goraptor.NewParser("guess")
	defer parserFile.Free()

	ch := parserFile.ParseFile(filename, "")
	for {
		statement, ok := <-ch
		if !ok {
			break
		}
		// fmt.Println(parser.Parse())
		// fmt.Printf("%#v\n", statement) // returning a goraptor.Statement (custom type)
		// fmt.Println(statement) //Basic Data structure?
		// abc := []byte{10}
		// fmt.Printf("%v", statement.GobDecode(abc))
		fmt.Printf("%#v", statement)
	}
}
