package sdprint

import (
	"fmt"
	"github.com/kr/pretty"
	"io"
)

func Pretty(v any) {
	fmt.Print(pretty.Sprint(v))
}

func PrettyL(v any) {
	fmt.Println(pretty.Sprint(v))
}

func FPretty(w io.Writer, v any) {
	_, _ = fmt.Fprint(w, pretty.Sprint(v))
}

func FPrettyL(w io.Writer, v any) {
	_, _ = fmt.Fprintln(w, pretty.Sprint(v))
}
