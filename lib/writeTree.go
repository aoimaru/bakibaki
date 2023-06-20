package lib

type Node struct {
	Path     string
	Children []*Node
}

type FileStatus struct {
	Path string
	Hash string
	Size uint32
	Mode uint32
}


func (index *Index)