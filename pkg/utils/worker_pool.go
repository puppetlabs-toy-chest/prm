package utils

import (
	"sync"
)

// Worker pool implementation adapted from https://brandur.org/go-worker-pool

// Pool is a worker group that runs a number of tasks at a
// configured concurrency.
type Pool[T any] struct {
	Tasks []*Task[T]

	concurrency int
	tasksChan   chan *Task[T]
	wg          sync.WaitGroup
}

func CreateWorkerPool[T any](tasks []*Task[T], workerCount int) *Pool[T] {
	if workerCount < 1 {
		workerCount = 1
	}

	return &Pool[T]{
		Tasks:       tasks,
		concurrency: workerCount,
		tasksChan:   make(chan *Task[T]),
	}
}

// Run runs all work within the pool and blocks until it's
// finished.
func (p *Pool[T]) Run() {
	for i := 0; i < p.concurrency; i++ {
		go p.work()
	}

	for _, task := range p.Tasks {
		p.wg.Add(1)
		p.tasksChan <- task
	}

	// all workers return
	close(p.tasksChan)

	p.wg.Wait()
}

// The work loop for any single goroutine.
func (p *Pool[T]) work() {
	for task := range p.tasksChan {
		task.Run(&p.wg)
	}
}

// Task encapsulates a work item that should go in a work
// pool.
type Task[T any] struct {
	// Output holds an error that occurred during a task. Its
	// result is only meaningful after Run has been called
	// for the pool that holds it.
	Name   string
	Output T

	f func() T
}

// Run runs a Task and does appropriate accounting via a
// given sync.WorkGroup.
func (t *Task[T]) Run(wg *sync.WaitGroup) {
	t.Output = t.f()
	wg.Done()
}

func CreateTask[T any](name string, f func() T, output T) *Task[T] {
	return &Task[T]{
		Name:   name,
		Output: output,
		f:      f,
	}
}
