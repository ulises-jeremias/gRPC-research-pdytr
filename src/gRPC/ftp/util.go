package ftp

import (
	"log"
	"os"
	"path/filepath"
)

// Djb2 function
func Djb2(bytes []byte) int64 {
	var hash int64 = 5381

	for _, c := range bytes {
		hash = ((hash << 5) + hash) + int64(c)
		// the above line is an optimized version of the following line:
		// hash = hash * 33 + int64(c)
		// which is easier to read and understand...
	}

	return hash
}

func FileHandler(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		*files = append(*files, path)
		return nil
	}
}
