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
	s, err := json.MarshalIndent(a, "", "    ")
	if err != nil {
		return fmt.Sprintf("RAW(%v)", a)
	}

	return string(s)
}
