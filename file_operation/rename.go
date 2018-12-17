package main

import (
	"log"
	"os"
)

func main() {
	/*文件移动和重命名*/
	originalPath := "../../test2.txt"
	newPath := "../../temp/test.txt"
	err := os.Rename(originalPath, newPath)
	if err != nil {
		log.Fatal(err)
	}
}
