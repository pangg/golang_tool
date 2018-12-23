package main

import (
	"fmt"
	"sync"
)

type Stack struct {
	stack []interface{}
	len   int
	lock  sync.Mutex
}

func main() {
	s := New()
	s.Push(1)
	s.Push("aaa")
	fmt.Println(s)

	fmt.Println(s.Pop())
	s.Push("bbbbbbbbbb")
	fmt.Println(s.len)

	fmt.Println(s.Peek())
}

func New() *Stack {
	s := &Stack{}
	s.stack = make([]interface{}, 0)
	s.len = 0
	return s
}

func (s *Stack) Len() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.len
}

func (s *Stack) isEmpty() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.len == 0
}

func (s *Stack) Pop() (el interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	el, s.stack = s.stack[0], s.stack[1:]
	s.len--
	return
}

func (s *Stack) Push(el interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	prepend := make([]interface{}, 1)
	prepend[0] = el
	s.stack = append(prepend, s.stack...)
	s.len++
}

func (s *Stack) Peek() interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.stack[0]
}
