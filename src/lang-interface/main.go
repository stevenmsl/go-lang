package main

/*
Interfaces specify behavior, and behavior is defined by set of methods
*/

import (
	"fmt"
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
func (BasicFileReadWriter) openTextFile(filename string) []byte {
	fmt.Println("File Read - BasicFileReadWriter")
	content := []byte{'d', 'o', 'n', 'e'}
	return content
}

func (BasicFileReadWriter) writeTextFile(filename string, content []byte) {
	fmt.Println("File Written - BasicFileReadWriter")
}

func main() {
	//rw is a variable of interface type FileReadWriter
	var rw FileReadWriter = BasicFileReadWriter{}
	rw.openTextFile("test.txt")
	rw.writeTextFile("test.txt", nil)
}
