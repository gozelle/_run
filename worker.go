package _run

type Worker interface {
	ID() int
	Run()
	Stop()
	Invoke(task Task)
}

type NewWorkerFunc func(id int, workerPool chan chan Task) Worker

func NewSimpleWorker(id int, workerPool chan chan Task) Worker {
	return &simpleWorker{
		id:          id,
		workerPool:  workerPool,
		taskChannel: make(chan Task),
		quit:        make(chan struct{}),
		closed:      make(chan struct{}),
	}
}

type simpleWorker struct {
	id          int
	taskChannel chan Task
	workerPool  chan<- chan Task
	quit        chan struct{}
	closed      chan struct{}
}

func (w *simpleWorker) ID() int {
	return w.id
}

func (w *simpleWorker) Run() {
	go func() {
		for {
			w.workerPool <- w.taskChannel

			select {
			case <-w.quit:
				w.closed <- struct{}{}
				return
			case task := <-w.taskChannel:
				w.Invoke(task)
			}
		}
	}()
}

func (w *simpleWorker) Invoke(task Task) {
	task.Invoke()
}

func (w *simpleWorker) Stop() {
	close(w.quit)
	<-w.closed
}
