package filesystem

import "os"

func GetFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
