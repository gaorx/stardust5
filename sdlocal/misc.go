package sdlocal

import (
	"os"
	"os/user"
	"runtime"
)

func Hostname() string {
	hn, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hn
}

func HomeDir() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	return usr.HomeDir
}

func OS() string {
	return runtime.GOOS
}

func Arch() string {
	return runtime.GOARCH
}

func NumCPU() int {
	return runtime.NumCPU()
}
