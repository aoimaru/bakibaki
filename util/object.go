package util

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetCommitObjectID() {
	current_dir, _ := os.Getwd()
	root := current_dir + "/.bakibaki/objects"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			rel_path, _ := filepath.Rel(root, path)
			fmt.Printf("path: %#v\n", rel_path)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
}
