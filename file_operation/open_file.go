package main

import (
	"log"
	"os"
)

func main() {
	//简单只读方式打开
	file, err := os.Open("../../temp/test.txt")
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	// OpenFile提供更多的选项。
	// 最后一个参数是权限模式permission mode
	// 第二个是打开时的属性
	file, err = os.OpenFile("../../temp/test.txt", os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	b := make([]byte, 1024)
	_, err = file.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(b))
	file.Close()
}
