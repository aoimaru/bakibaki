package lib

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	// "os"
	// "fmt"
	"strings"
)

func FindGitRoot(path string) (string, error) {
	filePaths, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}
	for _, filePath := range filePaths {
		if filePath.IsDir() && filePath.Name() == ".git" {
			absFilePath, err := filepath.Abs(filePath.Name())
			if err != nil {
				return "", err
			}
			return absFilePath, nil
		}
	}
	return "", errors.New("not git repository")
}

func FindBakiBakiRoot(path string) (string, error) {
	filePaths, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}
	for _, filePath := range filePaths {
		if filePath.IsDir() && filePath.Name() == ".bakibaki" {
			absFilePath, err := filepath.Abs(filePath.Name())
			if err != nil {
				return "", err
			}
			return absFilePath, nil
		}
	}
	return "", errors.New("not bakibaki repository")
}

func GetAllPath(root string) ([]string){
	files, err := ioutil.ReadDir(root)
    if err != nil {
        return nil
    }
	if strings.HasSuffix(root, "/.git") {
		return nil
	}
	if strings.HasSuffix(root, "/.bakibaki") {
		return nil
	}
	
	filePaths := make([]string, 0)

	for _, file := range files {
		if file.IsDir() {
			filePaths = append(filePaths, GetAllPath(filepath.Join(root, file.Name()))...)
			continue
		}
		filePaths = append(filePaths, filepath.Join(root, file.Name()))
	}
	return filePaths


}
