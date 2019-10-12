/*
package main
var res = []byte{
*/
package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"
	"testing"
)



func TestAAA(T *testing.T) {

	names := []string{"张三","离散"}
	namesStr:=strings.Join(names,",")
	fmt.Println(namesStr)


}

