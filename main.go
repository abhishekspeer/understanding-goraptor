package main

import ("fmt"
		"github.com/deltamobile/goraptor"
		"strconv"
)


func main(){
	defer fmt.Println("do something with statement")
	parser := goraptor.NewParser("guess")
	defer parser.Free()
	
	ch := parser.ParseUri("http://spdx.org/licenses/CC0-1.0", "")
	for {
    	statement, ok := <-ch
    	if ! ok {
        	break
    	}

	fmt.Println(statement)
	fmt.Println("-------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println(strconv.Itoa(parser.Free()))
	fmt.Println("-------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println(strconv.Itoa(goraptor.NewParser("guess")))


	
 	}

}