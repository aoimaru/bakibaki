package lib

import (
	"fmt"
	"os"
)

func CreateHEAD() error {
	current_dir, _ := os.Getwd()
	file_path := current_dir + "/.bakibaki/HEAD"
	write_object, err := os.Create(file_path)
	if err != nil {
		return err
	}
	defer write_object.Close()

	count, err := write_object.Write([]byte("ref: refs/heads/master\n"))
	if err != nil {
		return err
	}
	fmt.Printf("write %d bytes\n", count)
	return nil
}
