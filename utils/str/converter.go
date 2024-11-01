package str

import "encoding/json"

func DumpJSON(i interface{}) string {
	parsed, _ := json.Marshal(i)

	return string(parsed)
}
