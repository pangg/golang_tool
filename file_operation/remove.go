package main

import (
	"log"
	"os"
)

func main() {
	err := os.Remove("../../temp/test.txt")
	if err != nil {
		log.Fatal(err)
	}
}
