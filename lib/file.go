package lib

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
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

func (c *Client) CreateBlobFile(file_path string) ([]byte, string, error) {
	buffer, err := File2Byte(file_path, "blob")
	if err != nil {
		return nil, "", err
	}

	Pressed := Press(buffer)

	sha1 := sha1.New()
	sha1.Write(buffer)

	hash := hex.EncodeToString(sha1.Sum(nil))

	hashPath, err := hash2Path(hash)
	if err != nil {
		return nil, "", err
	}
	hashDir, err := hash2PathDir(hash)
	if err != nil {
		return nil, "", err
	}
	if _, err := os.Stat(c.Root + hashDir); err != nil {
		if err := os.MkdirAll(c.Root+hashDir, 1755); err != nil {
			return nil, "", err
		}
	}

	w, err := os.Create(c.Root + hashPath)
	if err != nil {
		return nil, "", err
	}
	defer w.Close()

	count, err := w.Write(Pressed)
	if err != nil {
		return nil, "", err
	}
	fmt.Println(hash)
	fmt.Printf("write %d bytes\n", count)

	return Pressed, hash, nil

}

func CreateTreeFile(file_path string) {

}

func CreateCommitFile(file_path string) {

}
