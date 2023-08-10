package sdtime

import (
	"time"
)

func SleepM(n int64) {
	time.Sleep(Minutes(n))
}

func SleepS(n int64) {
	time.Sleep(Seconds(n))
}

func SleepMs(n int64) {
	time.Sleep(Milliseconds(n))
}
