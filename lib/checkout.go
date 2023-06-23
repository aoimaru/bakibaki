package lib

import (
	"fmt"
	"io/ioutil"
	"os"
)

func (c *Client) GetHeadHash() string {
	ref_path, _ := GetHeadRef()
	// ref_path = "/mnt/c/Users/81701/Documents/AtCoder/.bakibaki/refs/heads/master"
	f, err := os.Open(ref_path)
	if err != nil {
		fmt.Println("ERROR:3", err)
		return ""
	}
	defer f.Close()
	buffer, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("ERROR:4", err)
		return ""
	}

	parent_hash := string(buffer)
	return parent_hash

}
