package sdtabledata

import (
	"encoding/json"
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdjson"
	"io/fs"
	"strings"
)

func trimDir(dir string) string {
	dir = strings.TrimPrefix(dir, "/")
	if dir == "" {
		dir = "."
	}
	return dir
}

func readJsonFile(fsys fs.FS, fn string, v any) error {
	jsonRaw, err := fs.ReadFile(fsys, fn)
	if err != nil {
		return sderr.Wrap(err, "read json data error")
	}
	jsonStr := strings.TrimSpace(string(jsonRaw))
	if jsonStr != "" {
		err := sdjson.UnmarshalString(jsonStr, v)
		if err != nil {
			return sderr.Wrap(err, "parse row data to json object error")
		}
	}
	return nil
}

func convertByJson(src, targetPtr any) error {
	rowsJson, err := json.Marshal(src)
	if err != nil {
		return sderr.Wrap(err, "rows data to json error")
	}
	err = json.Unmarshal(rowsJson, targetPtr)
	if err != nil {
		return sderr.Wrap(err, "unmarshal rows json error")
	}
	return nil
}
