
package lib

import (
	"errors"
	"io/ioutil"
	"path/filepath"
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