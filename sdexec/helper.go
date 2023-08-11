package sdexec

import (
	"bytes"
	"github.com/gaorx/stardust5/sdslog"
	"os/exec"
	"strconv"
	"time"

	"github.com/gaorx/stardust5/sderr"
)

func run(c *exec.Cmd, timeout time.Duration) error {
	if timeout <= 0 {
		err := c.Run()
		return sderr.Wrap(err, "run error")
	} else {
		err := c.Start()
		if err != nil {
			return sderr.Wrap(err, "start error")
		}
		done := make(chan error)
		go func() { done <- c.Wait() }()
		select {
		case <-time.After(timeout):
			err = c.Process.Kill()
			if err != nil {
				sdslog.WithError(err).With("pid", c.Process.Pid).Debug("kill timeout process error")
			}
			return ErrTimeout
		case err := <-done:
			return sderr.Wrap(err, "run error (with timeout)")
		}
	}
}

// 下面的代码主要复制于GO标准库，几乎没做修改

func combinedOutput(c *exec.Cmd, timeout time.Duration) ([]byte, error) {
	if c.Stdout != nil {
		return nil, sderr.New("stdout already set")
	}
	if c.Stderr != nil {
		return nil, sderr.New("stderr already set")
	}
	var b bytes.Buffer
	c.Stdout = &b
	c.Stderr = &b
	err := run(c, timeout)
	return b.Bytes(), err
}

func output(c *exec.Cmd, timeout time.Duration) ([]byte, error) {
	if c.Stdout != nil {
		return nil, sderr.New("stdout already set")
	}
	var stdout bytes.Buffer
	c.Stdout = &stdout

	captureErr := c.Stderr == nil
	if captureErr {
		c.Stderr = &prefixSuffixSaver{N: 32 << 10}
	}

	err := run(c, timeout)
	if err != nil && captureErr {
		if ee, ok := sderr.AsT[*exec.ExitError](err); ok {
			ee.Stderr = c.Stderr.(*prefixSuffixSaver).Bytes()
		}
	}
	return stdout.Bytes(), err
}

// prefixSuffixSaver is an io.Writer which retains the first N bytes
// and the last N bytes written to it. The Bytes() methods reconstructs
// it with a pretty error message.
type prefixSuffixSaver struct {
	N         int // max size of prefix or suffix
	prefix    []byte
	suffix    []byte // ring buffer once len(suffix) == N
	suffixOff int    // offset to write into suffix
	skipped   int64

	// TODO(bradfitz): we could keep one large []byte and use part of it for
	// the prefix, reserve space for the '... Omitting N bytes ...' message,
	// then the ring buffer suffix, and just rearrange the ring buffer
	// suffix when Bytes() is called, but it doesn't seem worth it for
	// now just for error messages. It's only ~64KB anyway.
}

func (w *prefixSuffixSaver) Write(p []byte) (n int, err error) {
	lenp := len(p)
	p = w.fill(&w.prefix, p)

	// Only keep the last w.N bytes of suffix data.
	if overage := len(p) - w.N; overage > 0 {
		p = p[overage:]
		w.skipped += int64(overage)
	}
	p = w.fill(&w.suffix, p)

	// w.suffix is full now if p is non-empty. Overwrite it in a circle.
	for len(p) > 0 { // 0, 1, or 2 iterations.
		n := copy(w.suffix[w.suffixOff:], p)
		p = p[n:]
		w.skipped += int64(n)
		w.suffixOff += n
		if w.suffixOff == w.N {
			w.suffixOff = 0
		}
	}
	return lenp, nil
}

// fill appends up to len(p) bytes of p to *dst, such that *dst does not
// grow larger than w.N. It returns the un-appended suffix of p.
func (w *prefixSuffixSaver) fill(dst *[]byte, p []byte) (pRemain []byte) {
	if remain := w.N - len(*dst); remain > 0 {
		add := min(len(p), remain)
		*dst = append(*dst, p[:add]...)
		p = p[add:]
	}
	return p
}

func (w *prefixSuffixSaver) Bytes() []byte {
	if w.suffix == nil {
		return w.prefix
	}
	if w.skipped == 0 {
		return append(w.prefix, w.suffix...)
	}
	var buf bytes.Buffer
	buf.Grow(len(w.prefix) + len(w.suffix) + 50)
	buf.Write(w.prefix)
	buf.WriteString("\n... omitting ")
	buf.WriteString(strconv.FormatInt(w.skipped, 10))
	buf.WriteString(" bytes ...\n")
	buf.Write(w.suffix[w.suffixOff:])
	buf.Write(w.suffix[:w.suffixOff])
	return buf.Bytes()
}
