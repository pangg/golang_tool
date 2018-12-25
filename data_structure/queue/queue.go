package queue

import (
	//"fmt"
	"sync"
)

type Queue struct {
	queue []interface{} //队列实际值
	len   int           //队列的长度
	lock  *sync.Mutex   //队列的锁
}

/*func main() {
	q := New()
	q.Push(1)
	q.Push(2)
	q.Push(3)
	fmt.Println(q)
	fmt.Println(q.Peek())
	fmt.Println(q.Shift())
	fmt.Println(q.Shift())
	fmt.Println(q.Shift())
}*/

func New() *Queue {
	queue := &Queue{}
	queue.queue = make([]interface{}, 0)
	queue.len = 0
	queue.lock = new(sync.Mutex)

	return queue
}

func (q *Queue) Len() int {
	//q.lock.Lock()
	//defer q.lock.Unlock()

	return q.len
}

func (q *Queue) isEmpty() bool {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.len == 0
}

func (q *Queue) Shift() (el interface{}) {
	el, q.queue = q.queue[0], q.queue[1:]
	q.len--
	return
}

func (q *Queue) Push(el interface{}) {
	q.queue = append(q.queue, el)
	q.len++
	return
}

func (q *Queue) Peek() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.queue[0]
}
