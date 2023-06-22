package lib

import (
	"fmt"
	"os"
)

func Test() {
	current_dir, _ := os.Getwd()
	ref_path := current_dir + "/.bakibaki/hello2.txt"
	w, err := os.Create(ref_path)
	if err != nil {
		fmt.Println(err)
	}
	defer w.Close()
}
