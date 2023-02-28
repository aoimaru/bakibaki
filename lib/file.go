
package lib

import (
	// "fmt"
	"bytes"
	"os"
	"strconv"
	"compress/zlib"
)

func Press(buffer []byte) []uint8 {
	var Pressed bytes.Buffer
	zWriter := zlib.NewWriter(&Pressed)
	zWriter.Write(buffer)
	zWriter.Close()

	return Pressed.Bytes()
}

func GetFileMeta(f *os.File, fType string) ([]byte, error){
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	size := strconv.FormatInt(info.Size(), 10)
	return []byte(fType+" "+size), nil
}


func File2Byte(file_path string, fType string) ([]byte, error) {
	f, err := os.Open(file_path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buffer := make([]byte, 0)
	buf := make([]byte, 64)
	for {
		n, _ := (*f).Read(buf)
		if n==0 {
			break
		}
		buffer = append(buffer, buf...)
	}
	meta, err := GetFileMeta(f, fType)
	if err != nil {
		return nil, err
	}
	buffer = append(meta, buffer...)

	Pressed := Press(buffer)
	return Pressed, nil
}


func CreateBlobFile(file_path string) ([]byte, error) {
	buffer, err := File2Byte(file_path, "blob")
	if err != nil {
		return nil, err
	}
	return buffer, nil

}

func CreateTreeFile(file_path string) {

}

func CreateCommitFile(file_path string) {

}

