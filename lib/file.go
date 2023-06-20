package lib

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"os"
	"strconv"
	// "io"
	// "reflect"
)

func Press(buffer []byte) []uint8 {
	var Pressed bytes.Buffer
	zWriter := zlib.NewWriter(&Pressed)
	zWriter.Write(buffer)
	zWriter.Close()

	return Pressed.Bytes()
}

func GetFileMeta(f *os.File, fType string) ([]byte, error) {
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	size := strconv.FormatInt(info.Size(), 10)
	return []byte(fType + " " + size), nil
}

func File2Byte(file_path string, fType string) ([]byte, error) {
	f, err := os.Open(file_path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	meta, err := GetFileMeta(f, fType)
	if err != nil {
		return nil, err
	}
	buffer, err := os.ReadFile(file_path)
	if err != nil {
		return nil, err
	}

	meta = append(meta, 0)
	buffer = append(meta, buffer...)
	return buffer, nil
}

func (c *Client_v2) CreateBlobFile(file_path string) {

}

func CreateTreeFile(file_path string) {

}

func CreateCommitFile(file_path string) {

}

func IsExist(file_path string) bool {
	if _, err := os.Stat(file_path); err == nil {
		return false
	}
	return true
}

func (gf *GitBuffer) ToFile(file_path string) error {
	write_object, err := os.Create(file_path)
	if err != nil {
		fmt.Println("OK:1")
		return err
	}
	defer write_object.Close()

	count, err := write_object.Write(gf.Buffer)
	if err != nil {
		fmt.Println("OK:2")
		return err
	}
	fmt.Printf("write %d bytes\n", count)
	return nil

}
