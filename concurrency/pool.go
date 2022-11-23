package concurrency

import (
	"sync"
)

type Task interface {
	Run(*sync.WaitGroup)
}

type Pool struct {
	tasks      []Task
	numThreads int
	tasksChan  chan Task
	wg         sync.WaitGroup
}

func NewPool(tasks []Task, numThreads int) *Pool {
	return &Pool{
		tasks:      tasks,
		numThreads: numThreads,
		tasksChan:  make(chan Task),
	}
}

func (pool *Pool) Run() {
	for i := 0; i < pool.numThreads; i++ {
		go pool.work()
	}

	pool.wg.Add(len(pool.tasks))
	for _, task := range pool.tasks {
		pool.tasksChan <- task
	}
	close(pool.tasksChan)

	pool.wg.Wait()
}

func (pool *Pool) work() {
	for task := range pool.tasksChan {
		task.Run(&pool.wg)
	}
}
