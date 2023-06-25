package lib

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
)

func WriteTree(node *Node, index *Index) string {
	if len((*node).Children) <= 0 {
		for _, entry := range *&index.Entries {
			if (*node).Path == entry.Name {
				return entry.Hash
			}
		}
	}

	buffer := make([]byte, 0)
	header := []byte{116, 114, 101, 101, 32, 51, 53, 51}
	buffer = append(buffer, header...)

	for _, child_node := range (*node).Children {
		file_status := child_node.GetFileStatus(index)
		entry_buffer := make([]byte, 0)
		entry_buffer = append(entry_buffer, 0)
		if file_status.GetType() == "blob" {
			entry_buffer = append(entry_buffer, []byte("blob"+" ")...)
			entry_buffer = append(entry_buffer, []byte(file_status.Name+" ")...)
			entry_buffer = append(entry_buffer, []byte(file_status.Hash)...)
		} else {
			entry_buffer = append(entry_buffer, []byte("tree"+" ")...)
			entry_buffer = append(entry_buffer, []byte(file_status.Name+" ")...)
			child_hash := WriteTree(child_node, index)
			entry_buffer = append(entry_buffer, []byte(child_hash)...)

		}
		buffer = append(buffer, entry_buffer...)
	}

	var compressed bytes.Buffer
	zlib_writer := zlib.NewWriter(&compressed)
	zlib_writer.Write(buffer)
	zlib_writer.Close()
	compressed_buffer := compressed.Bytes()

	sha1 := sha1.New()
	sha1.Write(compressed_buffer)

	new_hash := hex.EncodeToString(sha1.Sum(nil))

	current_dir, _ := os.Getwd()
	object_path := current_dir + "/.bakibaki/objects/"
	// fmt.Println(object_path, new_hash[:2], new_hash[2:])
	if _, err := os.Stat(object_path + new_hash[:2]); err != nil {
		if err := os.MkdirAll(object_path+new_hash[:2], os.ModePerm); err != nil {
			return ""
		}
	}

	if _, err := os.Stat(object_path + new_hash[:2] + "/" + new_hash[2:]); err == nil {
		return new_hash
	}

	new_writer, _ := os.Create(object_path + new_hash[:2] + "/" + new_hash[2:])
	defer new_writer.Close()

	count, err := new_writer.Write(compressed_buffer)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("write %d bytes\n", count)

	return new_hash
}

func CatFile(hash string) {
	current_dir, _ := os.Getwd()

	root_dir := current_dir + "/.bakibaki/objects/"
	tree_path := root_dir + hash[:2] + "/" + hash[2:]
	f, _ := os.Open(tree_path)
	defer f.Close()

	buffer := make([]byte, 0)
	buf := make([]byte, 64)
	for {
		n, _ := (*f).Read(buf)
		if n == 0 {
			break
		}
		buffer = append(buffer, buf...)
	}
	extracting_buffer := bytes.NewBuffer(buffer)
	zlib_f, _ := zlib.NewReader(extracting_buffer)

	zlib_buffer, _ := ioutil.ReadAll(zlib_f)
	entries := make([][]byte, 0)
	entry := make([]byte, 0)

	for _, zlib_buf := range zlib_buffer {
		if zlib_buf == 0 {
			entries = append(entries, entry)
			entry = make([]byte, 0)
		}
		entry = append(entry, zlib_buf)
	}
	entries = append(entries, entry)

	for _, entry := range entries {
		fmt.Println("entry:", string(entry))
	}
}
