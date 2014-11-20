// Copyright 2014, Ali Najafizadeh. All rights reserved.
// Use of this source code is governed by a BSD-style
// Author: Ali Najafizadeh
package note

import "fmt"

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func log(v ...interface{}) {
	fmt.Println(v)
}
