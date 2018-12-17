package main

import (
	"log"
	"os"
)

var (
	newFile *os.File
	err     error
)

func main() {
	newFile, err = os.Create("../../temp/test.txt")
	if err != nil {
		log.Fatal(err)
	}
	_, err := newFile.WriteString("hello words!!!!!!!!!!!!!!!!!!!")
	if err != nil {
		log.Println(err)
	}
	log.Println(newFile)
	newFile.Close()
}
