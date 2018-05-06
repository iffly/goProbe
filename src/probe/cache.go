package probe

import (
	"fmt"
	"sync"
	"time"
)

var (
	taskCache map[string][]Task
	taskMutex *sync.Mutex
)

func getTask(tp string, td time.Duration) (tasks []Task) {
	taskMutex.Lock()
	defer taskMutex.Unlock()

	_, ok := taskCache[tp]
	if !ok {
		fmt.Println("Failed to find *p")
		return
	}

	tasks = make([]Task, 0)

	for _, one := range taskCache[tp] {
		if td == one.invlReal {
			tasks = append(tasks, one)
		}
	}

	return
}

func setTask(tp string, tasks []Task, all bool) {
	taskMutex.Lock()
	defer taskMutex.Unlock()

	_, ok := taskCache[tp]
	if !ok {
		fmt.Println("Failed to find *p")
		return
	}

	if !all {
		for i, old := range taskCache[tp] {
			for _, new := range tasks {
				if old.key == new.key {
					taskCache[tp][i] = new
				}
			}
		}
		fmt.Println("tasks:", taskCache[tp])
	} else {
		for i, new := range tasks {
			for _, old := range taskCache[tp] {
				if old.key == new.key {
					tasks[i] = old
				}
			}
		}
		taskCache[tp] = tasks
		// fmt.Println("new task:", tasks)
	}
}
