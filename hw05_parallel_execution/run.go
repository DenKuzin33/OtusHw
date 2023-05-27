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
	wg.Add(len(tasks))
	workersPool := make(chan struct{}, n)
	maxErrors := int64(m)
	var errorsCount int64
	var result error

	for i, task := range tasks {
		if errorsCount >= maxErrors && maxErrors > 0 {
			// Вычтем незапущенные таски
			wg.Add(-(len(tasks) - i))
			result = ErrErrorsLimitExceeded
			break
		}
		workersPool <- struct{}{}
		task := task
		go func() {
			defer func() {
				<-workersPool
				wg.Done()
			}()
			taskError := task()
			if taskError != nil {
				atomic.AddInt64(&errorsCount, 1)
			}
		}()
	}

	wg.Wait()
	return result
}
