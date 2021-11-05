package _run

import (
	"time"
)

type Task interface {
	Invoke()
}

func NewDispatcher(taskQueueSize, queueSize int) *Dispatcher {
	return &Dispatcher{
		maxQueueSize: queueSize,
		taskQueue:    make(chan Task, taskQueueSize),
		workerPool:   make(chan chan Task, queueSize),
		workerRef:    make(map[int]Worker),
		quit:         make(chan struct{}),
		closed:       make(chan struct{}),
	}
}

type Dispatcher struct {
	maxQueueSize int
	taskQueue    chan Task
	workerPool   chan chan Task
	workerRef    map[int]Worker
	quit         chan struct{}
	closed       chan struct{}
}

func (d *Dispatcher) Send(task Task) {
	d.taskQueue <- task
}

func (d *Dispatcher) Run(fn NewWorkerFunc) {
	if fn == nil {
		fn = NewSimpleWorker
	}

	for i := 0; i < d.maxQueueSize; i++ {
		worker := fn(i+1, d.workerPool)
		worker.Run()
		d.workerRef[worker.ID()] = worker
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case task := <-d.taskQueue:
			go func(tsk Task) {
				workerCh := <-d.workerPool
				workerCh <- task
			}(task)
		case <-d.quit:
			for id := range d.workerRef {
				d.workerRef[id].Stop()
			}

			d.closed <- struct{}{}
			return
		}
	}
}

func (d *Dispatcher) Stop() {
	close(d.quit)
	select {
	case <-time.After(5 * time.Second):
		return
	case <-d.closed:
		return
	}
}
