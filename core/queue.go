package core

import "sync"

type Queue struct {
	Jobs   chan interface{} // Jobs is the channel that stores the items to be processed
	Errors chan error       // Errors is the channel that handle the errors
	task   QueueTaskRunner
	wg     sync.WaitGroup
}

type QueueTaskRunner interface {
	Run(job interface{}) error
}

func NewQueue(concurrent int, task QueueTaskRunner) *Queue {
	return &Queue{
		Jobs:   make(chan interface{}, concurrent),
		Errors: make(chan error),
		task:   task,
		wg:     sync.WaitGroup{},
	}
}

// Run executes the queue with n workers based in the concurrent number
func (q *Queue) Run() *Queue {
	for w := cap(q.Jobs); w > 0; w-- {
		q.wg.Add(1)
		go q.worker()
	}

	return q
}

// Wait all jobs to be processed
func (q *Queue) WaitWorkers() error {
	q.wg.Wait()
	if len(q.Errors) > 0 {
		return <-q.Errors
	}
	return nil
}

func (q *Queue) worker() {
	for {
		select {
		case job, ok := <-q.Jobs:
			if !ok {
				q.wg.Done()
				return
			}

			if err := q.task.Run(job); err != nil {
				q.wg.Done()
				q.Errors <- err
				return
			}
		}
	}
}
