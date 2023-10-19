package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	wg := new(sync.WaitGroup)
	var errsCount int32
	tasksChan := make(chan Task)

	for i := 0; i < n; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for task := range tasksChan {
				if err := task(); err != nil {
					atomic.AddInt32(&errsCount, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		if atomic.LoadInt32(&errsCount) >= int32(m) {
			break
		}
		tasksChan <- task
	}

	close(tasksChan)
	wg.Wait()

	if errsCount < int32(m) {
		return nil
	}

	return ErrErrorsLimitExceeded
}
