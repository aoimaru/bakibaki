package lib

import (
	"fmt"
	"os"
	"strings"
)

type Node struct {
	Path     string
	Children []*Node
}

type FileStatus struct {
	Name string
	Hash string
	Size uint32
	Mode uint32
}

func (node *Node) GetFileStatus(index *Index) FileStatus {
	current_dir, _ := os.Getwd()
	abs_path := strings.Replace(node.Path, "root", current_dir, 1)
	name := strings.Replace(node.Path, "root/", "", 1)
	var file_status FileStatus
	file_status.Name = name
	for _, entry := range index.Entries {
		if entry.Name == name {
			file_status.Hash = entry.Hash
			break
		}
	}

	f, _ := os.Open(abs_path)
	defer f.Close()

	if fi, err := f.Stat(); err == nil {
		file_status.Size = uint32(fi.Size())
		file_status.Mode = uint32(fi.Mode())
	}
	fmt.Printf("file_status: %+v\n", file_status)
	return file_status
}

func (fs *FileStatus) GetType() string {
	current_dir, _ := os.Getwd()
	file_path := current_dir + "/" + fs.Name
	file_buffer := make([]byte, 0)
	for _, file_buf := range []byte(file_path) {
		if file_buf == 0 {
			break
		}
		file_buffer = append(file_buffer, file_buf)
	}

	fmt.Println(string(file_buffer))

	if f, err := os.Stat(string(file_buffer)); os.IsNotExist(err) || f.IsDir() {
		return "tree"
	} else {
		return "blob"
	}

}

func GetParentName(tree *Node) string {
	tmp := strings.Split((*tree).Path, "/")
	return strings.Join(tmp[:len(tmp)-1], "/")
}

func (index *Index) CreateNodes() []*Node {
	var names []string

	current_dir, _ := os.Getwd()
	for _, entry := range index.Entries {
		file_path := current_dir + "/" + entry.Name
		if _, err := os.Stat(file_path); err != nil {
			continue
		}

		tmp := "root/" + entry.Name
		namespaces := strings.Split(tmp, "/")
		for i := 0; i <= len(namespaces); i++ {
			new_name := strings.Join(namespaces[:i], "/")
			flag := true
			for _, name := range names {
				if name == new_name {
					flag = false
					break
				}
			}
			if flag {
				names = append(names, new_name)
			}
		}
	}
	var nodes []*Node
	for _, name := range names {
		var node Node
		node.Path = name
		nodes = append(nodes, &node)
	}

	for _, node := range nodes {
		parent_path := GetParentName(node)

		for _, parent_node := range nodes {
			if (*parent_node).Path == parent_path {
				(*parent_node).Children = append((*parent_node).Children, node)
			}
		}
	}

	return nodes
}

type Column struct {
	Type string
	Name string
	Hash string
}

type Tree struct {
	Size    int
	Columns []Column
}

func (t *Tree) Format() {
	fmt.Printf("Object-Type: Tree  Size: %d\n", t.Size)
	for _, column := range t.Columns {
		fmt.Printf("%s %s %s\n", column.Type, column.Name, column.Hash)
	}
}

func (c *Client) GetTreeObject(hash string) Tree {
	// hash: byte列の先頭が0になってしまう。。。
	hash_buffer := []byte(hash)
	hash_string := hash
	if hash_buffer[0] == 0 {
		hash_string = string(hash_buffer[1:])
	}


	buffer, err := c.GetGitObject(hash_string)
	if err != nil {
		fmt.Println(err)
		return Tree{}
	}
	fmt.Println("OK", string(buffer))
	entry_buffer := make([][]byte, 0)
	entry_buf := make([]byte, 0)
	for _, buf := range buffer {
		if buf == 0 {
			entry_buffer = append(entry_buffer, entry_buf)
			entry_buf = make([]byte, 0)
		}
		entry_buf = append(entry_buf, buf)
	}
	entry_buffer = append(entry_buffer, entry_buf)

	for _, entry_buf := range entry_buffer[1:] {
		fmt.Println(string(entry_buf))
	}

	return Tree{}
}

func WalkingTree(hash string) {

}
