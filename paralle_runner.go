package _run

import "sync"

// NewParallelRunner
// 设置并行框架, parallelNum 设置并行数量，如果为0，默认50个并行数
// 设置等待数，到达等待数后停止wait
func NewParallelRunner(parallelNum, waitTotal int) *ParallelRunner {
	if parallelNum <= 0 {
		parallelNum = 50
	}
	
	p := &ParallelRunner{}
	p.barrier = make(chan struct{}, parallelNum)
	for i := 0; i < parallelNum; i++ {
		p.Release()
	}
	
	if waitTotal > 0 {
		var wg sync.WaitGroup
		p.wg = &wg
		p.Add(waitTotal)
	}
	
	return p
}

type ParallelRunner struct {
	barrier chan struct{}
	wg      *sync.WaitGroup
}

func (r *ParallelRunner) Acquire() {
	<-r.barrier
}

func (r *ParallelRunner) Release() {
	r.barrier <- struct{}{}
}

func (r *ParallelRunner) Run(lambda func()) {
	go func() {
		defer r.Done()
		r.Acquire()
		defer r.Release()
		lambda()
	}()
}

func (r *ParallelRunner) Wait() {
	if r.wg != nil {
		r.wg.Wait()
	}
}

func (r *ParallelRunner) Add(delta int) {
	if r.wg != nil {
		r.wg.Add(delta)
	}
}

func (r *ParallelRunner) Done() {
	if r.wg != nil {
		r.wg.Done()
	}
}
