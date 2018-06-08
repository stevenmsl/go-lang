/*

 */

package main

import (
	"fmt"
)

func main() {
	checkAssignability()
}

//T ...
type T [10]string

//S ...
type S [10]string

func acceptT(in T) {
	fmt.Println("In acceptT:")
	fmt.Println(in)
}
func accept(in [10]string) {
	fmt.Println("In accept:")
	fmt.Println(in)
}

func checkAssignability() {
	s := S{"a", "b", "c"}
	//acceptT(s) //Compilation error
	accept(s) //consider as same underlying types when at least one type is not named: S and [10]string (not named)
	t := T{"d", "e", "f"}
	acceptT(t)
	accept(t)
	//t = s //Compilation error
}
