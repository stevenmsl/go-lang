package main

import (
	"fmt"
)

func main() {
	growSlice()
}

func growSlice() {
	s := []string{"a", "b"} //a slice type has no specified length
	fmt.Printf("len %d\n", len(s))
	fmt.Printf("cap %d\n", cap(s))

	fmt.Println(s)
}
