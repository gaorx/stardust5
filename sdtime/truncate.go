package sdtime

import (
	"time"
)

func NowTruncateM() time.Time {
	return time.Now().Truncate(time.Minute)
}

func NowTruncateS() time.Time {
	return time.Now().Truncate(time.Second)
}

func NowTruncateMS() time.Time {
	return time.Now().Truncate(time.Millisecond)
}
