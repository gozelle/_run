package _run

import "time"

func NewForRunner(duration time.Duration) *ForRunner {
	return &ForRunner{ticker: time.NewTicker(duration)}
}

// ForRunner 模拟 for 循环执行，解决 cpu 百分百占用问题
// 默认按每秒触发
type ForRunner struct {
	ticker *time.Ticker
}

func (p *ForRunner) Stop() {
	if p.ticker != nil {
		p.ticker.Stop()
	}
}

func (p *ForRunner) Run(f func()) {
	if p.ticker == nil {
		p.ticker = time.NewTicker(time.Second)
	}
	for {
		f()
		select {
		case <-p.ticker.C:
			continue
		}
	}
}
