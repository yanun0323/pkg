package logs

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func getAbsPath(root string, dirs ...string) string {
	dir, _ := os.Getwd()
	path := dir

	if idx := strings.Index(dir, root); idx > 0 {
		path = dir[:strings.Index(dir, root)] + root
	}

	for _, dir := range dirs {
		path = fmt.Sprintf("%s/%s", path, dir)
	}

	return path
}

// Json formats the given object to a JSON string.
//
// If the object is not JSON serializable, it returns a string with the object's value.
func Json(a any) string {
	switch buf := a.(type) {
	case []byte:
		var data map[string]any
		_ = json.Unmarshal(buf, &data)
		if len(data) != 0 {
			a = data
		}
	case string:
		return Json([]byte(buf))
	case []rune:
		return Json([]byte(string(buf)))
	case byte:
		return Json([]byte{buf})
	case rune:
		return Json([]byte(string(buf)))
	default:
	}

	s, err := json.MarshalIndent(a, "", "    ")
	if err != nil {
		return fmt.Sprintf("RAW(%v)", a)
	}

	return string(s)
}
