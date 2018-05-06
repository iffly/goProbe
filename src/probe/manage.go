package probe

import (
	"fmt"
	"time"
)

type manage struct {
	taskChan chan Task
	retChan  chan Task
	thrdMax  int
	outChan  chan []Task
}

var (
	manageCache map[string]map[time.Duration]manage
	manageInvl  []time.Duration
)

func getProbeInfo(tp string, td time.Duration) (
	taskChan, retChan chan Task, outChan chan []Task, thrdMax int) {

	taskChan = manageCache[tp][td].taskChan
	retChan = manageCache[tp][td].retChan
	thrdMax = manageCache[tp][td].thrdMax

	outChan = manageCache[tp][td].outChan

	return
}

func adapterInvl(t time.Duration) (r time.Duration) {
	if t <= manageInvl[0] {
		r = manageInvl[0]
	}

	if t > manageInvl[len(manageInvl)-1] {
		r = manageInvl[len(manageInvl)-1]
	}

	for j := 1; j < len(manageInvl); j++ {
		if t > manageInvl[j-1] && t <= manageInvl[j] {
			r = manageInvl[j]
			break
		}
	}

	return
}

func ctorTask(tp string, tasks []Task) {
	for i := range tasks {
		if 5000 <= tasks[i].Rt {
			tasks[i].invlReal = manageInvl[len(manageInvl)-1]
		} else {
			tasks[i].invlReal = tasks[i].Invl
		}
	}

	setTask(tp, tasks, false)
}

func do(tp string, td time.Duration) {
	taskChan, retChan, outChan, thrdMax := getProbeInfo(tp, td)
	tasks := getTask(tp, td)
	if 0 == len(tasks) {
		return
	}

	inCnt := 0
	outCnt := 0

	for i := 0; i < len(tasks) && i < thrdMax; i++ {
		taskChan <- tasks[inCnt]
		inCnt++
	}

	rs := make([]Task, 0)
	for {
		if outCnt == len(tasks) {
			break
		}

		r := <-retChan
		outCnt++
		rs = append(rs, r)

		if inCnt < len(tasks) {
			taskChan <- tasks[inCnt]
			inCnt++
		}
	}

	ctorTask(tp, rs)

	outChan <- rs
	return
}

func manageThrd(tp string, td time.Duration) {
	for {
		tStart := time.Now()

		do(tp, td)

		if td <= time.Since(tStart) {
			fmt.Println("busy")
		} else {
			<-time.NewTimer(td - time.Since(tStart)).C
		}
	}
}

func UpdateTask(tp string, tasks []Task) {
	for i := range tasks {
		tasks[i].Invl = adapterInvl(tasks[i].Invl)
		tasks[i].invlReal = tasks[i].Invl
	}

	setTask(tp, tasks, true)
}
