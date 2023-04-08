package test

import (
	"fmt"
	"github.com/aoimaru/bakibaki/lib"
)


func WriteIndexHeaderTest(index *lib.Index) {
	bDirc := []byte((*index).Dirc)
	fmt.Println(bDirc)
}
