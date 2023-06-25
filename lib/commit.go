package lib

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aoimaru/bakibaki/util"
)

func GetCommitObject() {

}

func (c *Client) GetCommitObjectID() []string {
	current_dir, _ := os.Getwd()
	root := current_dir + "/.bakibaki/objects"

	hashes := make([]string, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			rel_path, _ := filepath.Rel(root, path)
			hash := strings.Replace(rel_path, "/", "", 1)
			buffer, _ := c.GetGitObject(hash)
			header := GetGitHeader(buffer)
			header = util.PaddingZeroBuffer(header)
			if strings.HasPrefix(header, "commit ") {
				hashes = append(hashes, hash)
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
	return hashes
}

func (c *Client) CreateCommitObject(buffer []byte) {
	entries := make([][]byte, 0)
	entry := make([]byte, 0)

	for _, buf := range buffer {
		if buf == 0 {
			entries = append(entries, entry)
			entry = make([]byte, 0)
		}
		entry = append(entry, buf)
	}
	entries = append(entries, entry)

	for _, entry := range entries {
		fmt.Println(string(entry))
	}
}

func (c *Client) GetCommitTreeHash(buffer []byte) (string, error) {
	entries := make([][]byte, 0)
	entry := make([]byte, 0)

	for _, buf := range buffer {
		if buf == 0 {
			entries = append(entries, entry)
			entry = make([]byte, 0)
		}
		entry = append(entry, buf)
	}
	entries = append(entries, entry)

	for _, entry := range entries {
		// fmt.Println(string(entry), entry)
		if strings.HasPrefix(string(entry), string([]byte{0, 116, 114, 101, 101})) {
			tree_hash := string(entry)
			tree_hash = strings.Replace(tree_hash, "tree ", "", 1)
			// fmt.Println("tree_hash:", tree_hash)
			return tree_hash, nil
		}
	}
	return "", errors.New("OK")
}
