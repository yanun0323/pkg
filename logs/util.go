package logs

import (
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
