package lexer

import "fmt"

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func log(v ...interface{}) {
	fmt.Println(v)
}
