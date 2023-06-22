package lib

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	// "io"
	// "reflect"
)

func IsExist(file_path string) bool {
	if _, err := os.Stat(file_path); err == nil {
		return false
	}
	return true
}

func (cf *IndexBuffer) ToFile(client Client) error {
	// indexファイルをオブジェクトとして取得
	index_path := client.GetIndexPath()
	write_object, err := os.Create(index_path)
	if err != nil {
		fmt.Println("OK:1")
		return err
	}
	defer write_object.Close()

	count, err := write_object.Write(cf.Buffer)
	if err != nil {
		fmt.Println("OK:2")
		return err
	}
	fmt.Printf("write %d bytes\n", count)
	return nil
}

func (cf *CommitBuffer) ToFile(client Client) (string, error) {
	buffer := cf.Buffer
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

	if _, err := os.Stat(object_path + new_hash[:2]); err != nil {
		if err := os.MkdirAll(object_path+new_hash[:2], 1755); err != nil {
			return "", err
		}
	}

	if _, err := os.Stat(object_path + new_hash[:2] + "/" + new_hash[2:]); err == nil {
		return "", err
	}

	new_writer, _ := os.Create(object_path + new_hash[:2] + "/" + new_hash[2:])
	defer new_writer.Close()

	count, err := new_writer.Write(compressed_buffer)
	if err != nil {
		return "", err
	}
	fmt.Printf("write %d bytes\n", count)

	return new_hash, nil

}
