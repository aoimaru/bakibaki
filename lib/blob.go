package lib

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
)

func (c *Client) CreateBlobFile(file_path string) ([]byte, string, error) {
	f, err := os.Open(file_path)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	// オブジェクトのヘッダーを生成するために, ファイルの情報を取得
	file_info, err := f.Stat()
	if err != nil {
		return nil, "", err
	}

	// ここでblobオブジェクトのバイト列を生成
	file_header := []byte("blob" + " " + strconv.FormatInt(file_info.Size(), 10))
	buffer, err := os.ReadFile(file_path)
	if err != nil {
		return nil, "", err
	}
	file_header = append(file_header, 0)
	buffer = append(file_header, buffer...)

	// Gitオブジェクトはzlibで圧縮されているので, 上記のblobオブジェクトのバイト列を圧縮して, 新しいバイト列を生成(compressed_buffer)
	var compressed bytes.Buffer
	zlib_writer := zlib.NewWriter(&compressed)
	zlib_writer.Write(buffer)
	zlib_writer.Close()
	compressed_buffer := compressed.Bytes()

	// sha1で, hash値を作成, 今回元はcompressed_bufferだけど, bufferと迷ってる(元のバイト列)
	sha1 := sha1.New()
	sha1.Write(compressed_buffer)
	hash := hex.EncodeToString(sha1.Sum(nil))

	// ハッシュ値から, ファイルパスと必要なディレクトリを作成する
	hash_rel_dir := "/objects/" + hash[:2]
	hash_rel_path := "/objects/" + hash[:2] + "/" + hash[2:]

	if _, err := os.Stat(c.Root + hash_rel_dir); err != nil {
		if err := os.MkdirAll(c.Root+hash_rel_dir, 1755); err != nil {
			return nil, "", err
		}
	}

	w, err := os.Create(c.Root + hash_rel_path)
	if err != nil {
		return nil, "", err
	}
	defer w.Close()

	count, err := w.Write(compressed_buffer)
	if err != nil {
		return nil, "", err
	}
	fmt.Println("hash:", hash)
	fmt.Printf("write: %d bytes\n", count)

	return compressed_buffer, hash, nil

}
