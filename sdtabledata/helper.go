package sdtabledata

import (
	"encoding/json"
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdjson"
	"os"
	"strings"
)

func readJsonFile(fnAbs string, v any) error {
	jsonRaw, err := os.ReadFile(fnAbs)
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
