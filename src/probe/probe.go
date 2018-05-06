package probe

import (
	"fmt"
	"time"
)

type Task struct {
	Invl     time.Duration
	invlReal time.Duration
	key      string
	Rt       int
}

var jfy = 0

func probe(task Task) (ret Task) {
	if 0 == jfy {
		task.Rt = 5000
	} else {
		task.Rt = 3000
	}
	jfy++

	tNow := time.Now()
	timeNow := tNow.Format("2006-01-02 15:04:05")
	fmt.Println(timeNow, "probe:", task.invlReal, task.key, task.Rt)
	ret = task
	return
}

func probeThrd(taskChan, retChan chan Task) {
	for {
		info := <-taskChan
		ret := probe(info)
		retChan <- ret
	}
}
