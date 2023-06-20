package util

// type Node struct {
// 	Path     string
// 	Name     string
// 	Parent   string
// 	Children []*Node
// }

// func ReversePath(levels []string) []string {
// 	reversed_levels := []string{}
// 	length := len(levels)
// 	for i := 1; i <= length; i++ {
// 		reversed_levels = append(reversed_levels, levels[length-i])
// 	}
// 	return reversed_levels
// }

// func DepthFirstWalking(node Node, count int) {
// 	fmt.Println()
// 	fmt.Println()
// 	// fmt.Printf("%+v\n", node)
// 	fmt.Println(count, "parent:", node.Name)
// 	fmt.Printf("%+v\n", node)
// 	for _, child := range node.Children {
// 		DepthFirstWalking(*child, count+1)
// 	}
// }

// func include(nodes []*Node, target *Node) bool {
// 	for _, node := range nodes {
// 		if (*node).Path == (*target).Path {
// 			return true
// 		}
// 	}
// 	return false
// }

// func CreateTree(index *lib.Index) {
// 	var nodes []*Node
// 	for _, entry := range (*index).Entries {
// 		levels := make([]string, 0)
// 		levels = append(levels, "root")
// 		levels = append(levels, strings.Split(entry.Name, "/")...)

// 		for i := 1; i < len(levels); i++ {
// 			new_node := Node{
// 				Path:   strings.Join(levels[:i+1], "/"),
// 				Name:   levels[i],
// 				Parent: strings.Join(levels[:i], "/"),
// 			}
// 			if !include(nodes, &new_node) {
// 				nodes = append(nodes, &new_node)
// 			}
// 		}
// 	}

// 	heads := make([]*Node, 0)
// 	for _, node := range nodes {
// 		fmt.Println()
// 		fmt.Printf("%+v\n", *node)
// 		for _, tmp := range nodes {
// 			if (*tmp).Path == (*node).Parent {
// 				fmt.Println("ddd:", *tmp)
// 				(*tmp).Children = append((*tmp).Children, node)
// 				fmt.Println("eee:", *tmp)
// 			}
// 		}
// 		if (*node).Parent == "root" {
// 			heads = append(heads, node)
// 		}
// 	}

// 	for _, head := range heads {
// 		DepthFirstWalking(*head, 0)
// 	}

// }
