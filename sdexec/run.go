package sdexec

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gaorx/stardust5/sderr"
)

type Result struct {
	Name   string            // 程序名称
	Args   []string          // 程序参数
	Dir    string            // 程序的运行目录
	Envs   map[string]string // 程序的环境变量
	Stdout []byte            // 标准输出
	Stderr []byte            // 错误输出
	Err    error             // 错误信息
	Usage  Rusage            // 进程概况
}

func (cmd *Cmd) New() *exec.Cmd {
	cmd1 := exec.Command(cmd.Name, cmd.Args...)
	if cmd.Dir != "" {
		cmd1.Dir = cmd.Dir
	}
	if len(cmd.Vars) > 0 {
		var env1 []string
		for name, val := range cmd.Vars {
			env1 = append(env1, fmt.Sprintf("%s=%s", name, val))
		}
		cmd1.Env = env1
	}
	return cmd1
}

var (
	ErrTimeout = sderr.New("timeout")
)

func (cmd *Cmd) Run() error {
	cmd1 := cmd.New()
	return run(cmd1, cmd.Timeout)
}

func (cmd *Cmd) RunConsole() error {
	cmd1 := cmd.New()
	cmd1.Stdin = os.Stdin
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr
	return run(cmd1, cmd.Timeout)
}

func (cmd *Cmd) RunResult() *Result {
	cmd1 := cmd.New()
	stdout, stderr := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	cmd1.Stdout, cmd1.Stderr = stdout, stderr
	err := run(cmd1, cmd.Timeout)
	var r Result
	r.Name = cmd.Name
	r.Args = cmd.Args
	r.Dir = cmd.Dir
	r.Envs = cmd.Vars
	r.Stdout = stdout.Bytes()
	r.Stderr = stderr.Bytes()
	r.Err = err
	if cmd1.ProcessState != nil {
		r.Usage = getRusage(cmd1.ProcessState.SysUsage())
	} else {
		r.Usage = Rusage{}
	}
	return &r
}

func (cmd *Cmd) RunOutput(combine bool) ([]byte, error) {
	cmd1 := cmd.New()
	var out []byte
	var err error
	if combine {
		out, err = combinedOutput(cmd1, cmd.Timeout)
	} else {
		out, err = output(cmd1, cmd.Timeout)
	}
	return out, err
}

func (cmd *Cmd) RunOutputString(combine bool) (string, error) {
	buff, err := cmd.RunOutput(combine)
	if err != nil {
		return "", err
	}
	return string(buff), nil
}

func (r *Result) HasErr() bool {
	return r.Err != nil
}

func (r *Result) StdoutString() string {
	return string(r.Stdout)
}

func (r *Result) StderrString() string {
	return string(r.Stderr)
}

func (r *Result) StdoutLines() []string {
	return strings.Split(r.StdoutString(), "\n")
}

func (r *Result) StderrLines() []string {
	return strings.Split(r.StderrString(), "\n")
}

func (r *Result) ExitCode() int {
	if r.Err == nil {
		return 0
	}
	if exitErr, ok := r.Err.(*exec.ExitError); ok {
		return exitErr.ExitCode()
	} else {
		return -9999
	}
}

func (r *Result) Cli() string {
	cliEscape := func(s string) string {
		if strings.Contains(s, " ") {
			return "\"" + s + "\""
		} else {
			return s
		}
	}
	buf := bytes.NewBufferString("")
	buf.WriteString(cliEscape(r.Name))
	for _, arg := range r.Args {
		buf.WriteString(" ")
		buf.WriteString(cliEscape(arg))
	}
	return buf.String()
}
