package content

import (
	"io/fs"
	"os"
)

var (
	fileSystem fs.FS
)

func init() {
	fileSystem = os.DirFS(".")
}

// FS returns the content filesystem
func FS() fs.FS {
	return fileSystem
}
