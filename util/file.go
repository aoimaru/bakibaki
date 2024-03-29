package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetGitObjectHeader(buffer *[]byte) (string, error) {

	header := make([]byte, 0)
	for _, buf := range *buffer {
		if buf == 0 {
			break
		}
		header = append(header, buf)
	}
	return string(header), nil

}

func WalkingDir(root string) ([]string, error) {
	if _, err := os.Stat(root); err != nil {
		return nil, err
	}
	paths := make([]string, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		rel_path, err := filepath.Rel(root, path)
		if info.IsDir() {
			if strings.HasPrefix(rel_path, ".bakibaki") {
				return filepath.SkipDir
			}
			return nil
		}
		paths = append(paths, rel_path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return paths, nil
}

func GetBrabches(root string) ([]string, error) {
	if _, err := os.Stat(root); err != nil {
		fmt.Println(err)
	}
	branches := make([]string, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		rel_path, err := filepath.Rel(root, path)
		if rel_path != "." {
			// fmt.Println("rel_path:", rel_path)
			branches = append(branches, rel_path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return branches, nil
}

func PaddingZeroBuffer(tag string) string {
	buffer := make([]byte, 0)
	for _, buf := range []byte(tag) {
		if buf != 0 {
			buffer = append(buffer, buf)
		}
	}
	return string(buffer)
}
