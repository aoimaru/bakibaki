package lib

import (
	"fmt"
	"os"
)

func (c *Client) UpdateRef(head string, hash string) error {
	current_dir, _ := os.Getwd()
	head_path := current_dir + "/.bakibaki/" + head
	if _, err := os.Stat(head_path); err != nil {
		return err
	}

	write_object, err := os.Create(head_path)
	if err != nil {
		fmt.Println("ERROR:1", err)
		return err
	}
	defer write_object.Close()

	if _, err = write_object.Write([]byte(hash)); err != nil {
		fmt.Println("ERROR:2", err)
		return err
	}

	return nil
}
