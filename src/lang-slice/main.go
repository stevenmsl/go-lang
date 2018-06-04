package main

/*
A slice is a descriptor of an array segment. It consists of a pointer to the array, the length of the segment,
and its capacity (the maximum length of the segment).
*/

import (
	"fmt"
)

func main() {
	growSlice()
}

func growSlice() {
	s := []string{"a", "b"} //a slice type has no specified length
	fmt.Println("before growing:")
	fmt.Printf("len %d\n", len(s))
	fmt.Printf("cap %d\n", cap(s))
	fmt.Println(s)
	var _cap int
	if _cap = cap(s); _cap == 0 {
		_cap = 1
	}
	t := make([]string, len(s), _cap*2)
	/*
		The built-in copy function supports copying between slices of different lengths (it will copy only up to the smaller number of elements).
		In addition, copy can handle source and destination slices that share the same underlying array, handling overlapping slices correctly.
	*/
	copy(t, s)
	s = t
	fmt.Println("after growing:")
	fmt.Printf("len %d\n", len(s))
	fmt.Printf("cap %d\n", cap(s))
	fmt.Println(s)
}
