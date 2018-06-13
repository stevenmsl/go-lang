package main

import (
	"fmt"
)

//FileReader ...
type FileReader interface {
	openTextFile(filename string) []byte
}

//StringFileReader ...
type StringFileReader interface {
	readToEnd(filename string) string
	FileReader
}

//XFileReader ...
type XFileReader interface {
	readToEnd(filename string) string
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

func (BasicFileReadWriter) openTextFile(filename string) []byte {
	fmt.Println("File Read - BasicFileReadWriter")
	content := []byte{'d', 'o', 'n', 'e'}
	return content
}

func (BasicFileReadWriter) writeTextFile(filename string, content []byte) {
	fmt.Println("File Written - BasicFileReadWriter")
}

//BasicXFileReader ...
type BasicXFileReader struct {
}

func (BasicXFileReader) openTextFile(filename string) []byte {
	fmt.Println("File Read - BasicFileReadWriter")
	content := []byte{'d', 'o', 'n', 'e'}
	return content
}

func (BasicXFileReader) readToEnd(filename string) string {
	fmt.Println("File Read - BasicXFileReader")
	content := "done"
	return content
}

//BasicStringFileReader ...
type BasicStringFileReader struct {
	format string
}

func (BasicStringFileReader) openTextFile(filename string) []byte {
	fmt.Println("File Read - BasicFileReadWriter")
	content := []byte{'d', 'o', 'n', 'e'}
	return content
}

func (BasicStringFileReader) readToEnd(filename string) string {
	fmt.Println("File Read - BasicXFileReader")
	content := "done"
	return content
}

//BasicFileReader ...
type BasicFileReader struct {
}

func (BasicFileReader) openTextFile(filename string) []byte {
	fmt.Println("File Read - BasicFileReadWriter")
	content := []byte{'d', 'o', 'n', 'e'}
	return content
}

func main() {
	typeAssertionConcrete()
	typeAssertionInterface()

}

//the 3rd assignability case
func checkAssignabilityCase3() {
	var rw FileReadWriter = BasicFileReadWriter{}
	var r FileReader = rw //FileReader is a subset of FileReadWriter. So this is fine.
	_ = r                 //to bypass declared but not used error
}

func typeConversionReader() {
	var xr XFileReader = BasicXFileReader{}
	//You can use a XFileReader as a StringFileReader
	//as both interfaces have the have the same method set
	var sr StringFileReader = xr
	//You can use a StringFileReader as a FileReader
	//as FileReader is a subset of StringFileReader
	var r FileReader = sr
	_ = r
	r = BasicFileReader{}
	//This is not allowed. You cannot treat a FileReader as a StringFileReader as the readToEnd method is missing
	//from the FileReader interface.
	//sr = r
}

func typeAssertionConcrete() {
	fmt.Println("In typeAssertionConcrete...")
	var r FileReader = BasicFileReadWriter{}
	//The following will fail the type assertion
	//as complier wouldn’t know if r can be considered as a BasicXFileReader -
	//BasicXFileReader can implement methods that are not defined in the FileReader interface.
	//var xr BasicXFileReader = r

	//type assertion expression
	c := r.(BasicFileReadWriter)
	fmt.Printf("Interface: FileReader, Concrete:%T\n", c)
	//a BasicFileReadWriter is not a BasicFileReader - both types are concrete
	br, ok := r.(BasicFileReader)
	if !ok {
		fmt.Printf("ok: %v\n", ok)
		fmt.Printf("%v, %T\n", br, br)
	}

}

func typeAssertionInterface() {
	fmt.Println("In typeAssertionInterface...")
	var r FileReader = BasicStringFileReader{"json"}
	var sr StringFileReader
	//It checks if the dynamic value (BasicStringFileReader, concrete type) satisfies desired interface
	//and returns value of such interface type value (StringFileReader, interface type).
	sr, ok := r.(StringFileReader)
	fmt.Printf("%T %v %v\n", sr, sr, ok)

	var rw FileReadWriter
	//BasicStringFileReader does not satisfy FileReadWriter interface
	rw, ok = r.(FileReadWriter)
	fmt.Printf("%T %v %v\n", rw, rw, ok)

}
