package lib

import (
	"fmt"
	"os"
)

func (c *Client) UpdateRef(head Head, hash string) error {
	// current_dir, _ := os.Getwd()
	var head_path string
	switch v := head.(type) {
	case DetachedHead:
		head_path = "/.bakibaki/HEAD"
	case TatchedHead:
		head_path = v.Head
	}
	head_buffer := make([]byte, 0)
	for _, head_buf := range []byte(head_path) {
		if head_buf == 0 {
			break
		}
		head_buffer = append(head_buffer, head_buf)
	}
	head_path = string(head_buffer)

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
