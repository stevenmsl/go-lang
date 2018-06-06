package main

/*
A slice is a descriptor of an array segment. It consists of a pointer to the array, the length of the segment,
and its capacity (the maximum length of the segment).
*/

import (
	"fmt"
)

func main() {
	s := []string{"a", "b", "c"}
	fmt.Println("Before appending...")
	fmt.Println(s)
	s = append(s, "d", "e")
	fmt.Println("After appending...")
	fmt.Println(s)

	slice()
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

func slice() {
	/*
		Slicing does not copy the slice's data.
		It creates a new slice value that points to the original array.
		This makes slice operations as efficient as manipulating array indices.
		Therefore, modifying the elements (not the slice itself) of a re-slice modifies the elements of the original slice:
	*/
	s := []string{"a", "b", "c", "d", "e"}
	t := s[2:]
	fmt.Printf("len(s):%d\n", len(s))
	fmt.Printf("cap(s):%d\n", cap(s))
	fmt.Print("s:")
	fmt.Println(s)
	fmt.Printf("len(t):%d\n", len(t))
	fmt.Printf("cap(t):%d\n", cap(t))
	fmt.Print("t:")
	fmt.Println(t)
	t[0] = "c1"
	fmt.Println("After s is modified ")
	fmt.Print("s:")
	fmt.Println(s)
}

func append(slice []string, data ...string) []string {
	m := len(slice)
	n := m + len(data)
	if n > cap(slice) { //Double the size if needed
		newSlice := make([]string, n*2)
		copy(newSlice, slice)
		slice = newSlice
	}
	copy(slice[m:n], data)
	return slice
}
