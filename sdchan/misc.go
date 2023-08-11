package sdchan

import (
	"reflect"
)

func MergeRecv[T any](chans ...<-chan T) <-chan T {
	removeChan := func(chans []<-chan T, index int) []<-chan T {
		if len(chans) <= 0 {
			return nil
		}
		r := make([]<-chan T, 0)
		for i, c := range chans {
			if i != index {
				r = append(r, c)
			}
		}
		return r
	}
	if len(chans) <= 0 {
		return nil
	}
	r := make(chan T)
	go func() {
		defer close(r)
		chans1 := removeChan(chans, -1) // Copy
		for {
			chosen, v, recvOK := SelectOne(chans1...)
			if !recvOK {
				chans1 = removeChan(chans1, chosen)
				if len(chans1) <= 0 {
					break
				}
			} else {
				r <- v
			}
		}
	}()
	return r
}

func SelectOne[T any](chans ...<-chan T) (int, T, bool) {
	if len(chans) <= 0 {
		var empty T
		return -1, empty, false
	}
	cl := make([]reflect.SelectCase, 0, len(chans))
	for _, c := range chans {
		cl = append(cl, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(c),
		})
	}
	index, v, ok := reflect.Select(cl)
	return index, v.Interface().(T), ok
}
