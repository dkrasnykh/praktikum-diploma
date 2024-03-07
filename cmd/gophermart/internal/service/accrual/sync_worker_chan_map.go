package accrual

import "sync"

type WorkerTimeoutChan struct {
	value map[int]chan int
	mx    sync.RWMutex
}

func NewWorkerChanMap(rateLimit int) *WorkerTimeoutChan {
	return &WorkerTimeoutChan{
		value: make(map[int]chan int, rateLimit),
		mx:    sync.RWMutex{},
	}
}

func (m *WorkerTimeoutChan) Insert(key int, value chan int) {
	m.mx.Lock()
	defer m.mx.Unlock()

	m.value[key] = value
}

func (m *WorkerTimeoutChan) Get(key int) chan int {
	m.mx.RLock()
	defer m.mx.RUnlock()

	return m.value[key]
}

func (m *WorkerTimeoutChan) Broadcast(timeout int) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	for _, ch := range m.value {
		ch <- timeout
	}
}
