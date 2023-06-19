package lib

import (
	"bytes"
	"compress/zlib"
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
