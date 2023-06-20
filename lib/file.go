package lib

import (
	"fmt"
	"os"
	// "io"
	// "reflect"
)

func IsExist(file_path string) bool {
	if _, err := os.Stat(file_path); err == nil {
		return false
	}
	return true
}

func (gf *GitBuffer) ToFile(file_path string) error {
	write_object, err := os.Create(file_path)
	if err != nil {
		fmt.Println("OK:1")
		return err
	}
	defer write_object.Close()

	count, err := write_object.Write(gf.Buffer)
	if err != nil {
		fmt.Println("OK:2")
		return err
	}
	fmt.Printf("write %d bytes\n", count)
	return nil

}
