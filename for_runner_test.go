package _run

import (
	"log"
	"testing"
	"time"
)

func TestFor_Run(t *testing.T) {
	fr := NewForRunner(1 * time.Second)
	go fr.Run(func() {
		log.Println("ok1")
	})
	go fr.Run(func() {
		log.Println("ok2")
		time.Sleep(5 * time.Second)
	})
	time.Sleep(20 * time.Second)
	t.Log("准备停止运行...")
	fr.Stop()
	t.Log("已停止运行...")
}
