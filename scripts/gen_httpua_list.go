package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	resp, err := http.Get("https://www.useragentstring.com/pages/useragentstring.php?name=All")
	if err != nil {
		panic(err)
	}

	htmlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	html := string(htmlBytes)

	patt := regexp.MustCompile(`<a href='/(\w|\.)+_id_\d+.php'>([^<]+)</a>`)
	l := patt.FindAllStringSubmatch(html, -1)
	var lines []string
	for _, ss := range l {
		ua := ss[2]
		lines = append(lines, fmt.Sprintf(`		%s,`, strconv.Quote(ua)))
	}

	t := `
package sdhttpua
var (
	rawUserAgents = []string{
%s
	}
)
`
	goFile := fmt.Sprintf(t, strings.Join(lines, "\n"))
	fmt.Println(goFile)
}
