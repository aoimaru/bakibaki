package lib

import (
	"fmt"
	"os"
)

func (c *Client) GetHeadHash() string {
	// ref_path, _ := GetHeadRef()
	// // ref_path = "/mnt/c/Users/81701/Documents/AtCoder/.bakibaki/refs/heads/master"
	// f, err := os.Open(ref_path)
	// if err != nil {
	// 	fmt.Println("ERROR:3", err)
	// 	return ""
	// }
	// defer f.Close()
	// buffer, err := ioutil.ReadAll(f)
	// if err != nil {
	// 	fmt.Println("ERROR:4", err)
	// 	return ""
	// }

	// parent_hash := string(buffer)
	// return parent_hash

	current, _ := os.Getwd()
	BakiBakiRootPath, err := FindBakiBakiRoot(current)
	if err != nil {
		fmt.Println(err)
	}

	client := Client{
		Root: BakiBakiRootPath,
	}
	commit_hash, err := CommitTree("", "dnfpcauds", client)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("commit hash:", commit_hash)
	return commit_hash

}
