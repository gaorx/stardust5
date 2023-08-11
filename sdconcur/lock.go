package sdconcur

import (
	"sync"
)

func Lock(mtx *sync.Mutex, action func()) {
	if mtx != nil {
		mtx.Lock()
		defer mtx.Unlock()
	}
	action()
}

func LockW(mtx *sync.RWMutex, action func()) {
	if mtx != nil {
		mtx.Lock()
		defer mtx.Unlock()
	}
	action()
}

func LockR(mtx *sync.RWMutex, action func()) {
	if mtx != nil {
		mtx.RLock()
		defer mtx.RUnlock()
	}
	action()
}
