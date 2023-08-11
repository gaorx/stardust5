package sdprint

import (
	"fmt"
	"github.com/gaorx/stardust5/sdjson"
	"io"
)

func Json(v any, pretty bool) {
	if pretty {
		fmt.Print(sdjson.MarshalPretty(v))
	} else {
		fmt.Print(sdjson.MarshalStringDef(v, ""))
	}
}

func JsonL(v any, pretty bool) {
	if pretty {
		fmt.Println(sdjson.MarshalPretty(v))
	} else {
		fmt.Println(sdjson.MarshalStringDef(v, ""))
	}
}

func FJson(w io.Writer, v any, pretty bool) (int, error) {
	if pretty {
		return fmt.Fprint(w, sdjson.MarshalPretty(v))
	} else {
		return fmt.Fprint(w, sdjson.MarshalStringDef(v, ""))
	}
}

func FJsonL(w io.Writer, v any, pretty bool) (int, error) {
	if pretty {
		return fmt.Fprintln(w, sdjson.MarshalPretty(v))
	} else {
		return fmt.Fprintln(w, sdjson.MarshalStringDef(v, ""))
	}
}
