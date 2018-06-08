package main

/*
Interfaces specify behavior, and behavior is defined by set of methods.
*/

/*
In Go, types which implement all interface’s method automatically satisfy such interface.
You don’t need to use keywords such as implements to say that certain type implements a certain interface.
*/

import (
	"fmt"
	"reflect"
)

//FileReader ...
type FileReader interface {
	openTextFile(filename string) []byte
}

//FileWriter ...
type FileWriter interface {
	writeTextFile(filename string, content []byte)
}

//FileReadWriter ...
type FileReadWriter interface {
	//composition of interfaces
	FileReader
	FileWriter
}

//BasicFileReadWriter ...
type BasicFileReadWriter struct {
}

//Add functions to the struct to implement the interfaces
//You don’t need to use keywords such as implements to say that certain type implements a certain interface.
func (BasicFileReadWriter) openTextFile(filename string) []byte {
	fmt.Println("File Read - BasicFileReadWriter")
	content := []byte{'d', 'o', 'n', 'e'}
	return content
}

func (BasicFileReadWriter) writeTextFile(filename string, content []byte) {
	fmt.Println("File Written - BasicFileReadWriter")
}

func checkDynamicType(i interface{}) {
	dt := reflect.TypeOf(i)
	fmt.Printf("PkgPath: %s, Name: %s, Type: %s\n", dt.PkgPath(), dt.Name(), dt.String())
}

func checkNil() {
	rw := returnInterface()
	dt := reflect.TypeOf(rw)
	fmt.Printf("PkgPath: %s, Name: %s, Type: %s\n", dt.PkgPath(), dt.Name(), dt.String())

	if rw == nil {
		fmt.Println("returnInterface returns: rw is nil")
	} else {
		fmt.Println("returnInterface returns: rw is not nil")
	}
}

//this function returns interface FileReadWriter
func returnInterface() FileReadWriter {
	var t *BasicFileReadWriter //not initialized

	/*
		if false {
			t = &BasicFileReadWriter{}
		}
	*/
	if t == nil {
		fmt.Println("Inside returnInterface: t is null")
	}

	return t //return interface value which is not nil but holds a nil pointer
}

func main() {
	checkNil()
	//rw is a variable of interface type FileReadWriter
	var rw FileReadWriter = BasicFileReadWriter{}
	checkDynamicType(rw)
	rw.openTextFile("test.txt")
	rw.writeTextFile("test.txt", nil)

}
