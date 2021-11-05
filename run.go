package _run

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

var (
	names = []string{
		"三国杀",
		"阿莫西林",
		"杜鹃",
		"梁朝伟",
		"康熙字典",
		"龙舟赛",
		"风扇",
		"洗手液",
		"肥皂",
		"电影",
		"运动",
		"游泳",
		"比赛",
		"大国庆典",
		"台湾岛",
		"果冻",
		"美女",
		"下饭",
		"蛇仙姐姐",
		"手机号",
		"叮咚买菜",
		"年度报告",
		"2025",
		"007",
		"茶叶",
		"香蕉",
		"米饭",
		"庄周",
		"马可波罗",
		"张飞",
	}
)

func main() {
	d := NewDispatcher(10, 5)
	defer d.Stop()
	d.Run(nil)
	
	for i := 0; i < 30; i++ {
		go func(idx int) {
			d.Send(&simpleTask{
				seq:  idx + 1,
				name: names[idx],
			})
		}(i)
	}
	
	s := make(chan os.Signal)
	signal.Notify(s, os.Interrupt, os.Kill)
	<-s
}

type simpleTask struct {
	seq  int
	name string
}

func (s *simpleTask) Invoke() {
	rand.Seed(time.Now().UnixNano())
	sec := rand.Intn(50)
	time.Sleep(time.Duration(sec) * time.Millisecond)
	fmt.Printf("[%02d]: %s\n", s.seq, s.name)
}
