package probe

import (
	"sync"
	"time"
)

func Init() {
	taskCache = make(map[string][]Task, 0)
	taskMutex = new(sync.Mutex)

	manageCache = make(map[string]map[time.Duration]manage)
	manageInvl = []time.Duration{
		5 * time.Second,
		10 * time.Second,
		20 * time.Second,
		30 * time.Second,
		60 * time.Second,
	}
}

func Reg(tp string) (outChan chan []Task) {
	taskMutex.Lock()
	defer taskMutex.Unlock()

	_, ok := taskCache[tp]
	if ok {
		panic("double reg")
	}
	taskCache[tp] = make([]Task, 0)

	_, ok = manageCache[tp]
	if !ok {
		manageCache[tp] = make(map[time.Duration]manage)
	}

	outChan = make(chan []Task)
	for _, invl := range manageInvl {
		manageCache[tp][invl] = manage{
			taskChan: make(chan Task),
			retChan:  make(chan Task),
			// thrdMax:  int(1024 * time.Second / invl),
			thrdMax: 1,
			outChan: outChan,
		}
		for i := 0; i < manageCache[tp][invl].thrdMax; i++ {
			go probeThrd(manageCache[tp][invl].taskChan, manageCache[tp][invl].retChan)
		}
		go manageThrd(tp, invl)
	}
	return outChan
}

func Main() {
	Init()
	outChan := Reg("udpProbe")

	go func() {
		t1 := Task{
			Invl: 4 * time.Second,
			key:  "t1",
		}
		UpdateTask("udpProbe", []Task{t1})
	}()

	for {
		<-outChan
		// info := <-outChan
		// fmt.Println("ret:", info)
	}
}
