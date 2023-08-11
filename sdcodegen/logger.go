package sdcodegen

import (
	"fmt"
	"github.com/gaorx/stardust5/sdslog"
)

func Slog(action Action, fn, absFn string) {
	switch action {
	case I:
		sdslog.Infof("IGNORE %s", fn)
	case W:
		sdslog.Infof("WRITE  %s", fn)
	case C:
		sdslog.Infof("CREATE %s", fn)
	}
}

func SimplePrint(action Action, fn, absFn string) {
	switch action {
	case I:
		fmt.Printf("IGNORE %s\n", fn)
	case W:
		fmt.Printf("WRITE  %s\n", fn)
	case C:
		fmt.Printf("CREATE %s\n", fn)
	}
}
