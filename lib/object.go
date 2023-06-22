package lib

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	LENGTH_OF_HASH       = 2
	LENGTH_OF_HEADER     = 1
	LENGTH_OF_DATA       = 2
	NUM_OF_COMMIT_COLUMN = 2
)

var (
	emailRegexpString     = "([a-zA-Z0-9_.+-]+@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\\.)+[a-zA-Z]{2,})"
	timestampRegexpString = "([1-9][0-9]* \\+[0-9]{4})"
	sha1Regexp            = regexp.MustCompile("[0-9a-f]{20}")
	signRegexp            = regexp.MustCompile("^[^<]* <" + emailRegexpString + "> " + timestampRegexpString + "$")
)

type GitBuffer struct {
	Buffer []byte
}

type IndexBuffer struct {
	Buffer []byte
}

type CommitBuffer struct {
	Buffer []byte
}

type GitObject interface {
	Format()
}

type Column struct {
	Type string
	Name string
	Hash string
}

type Tree struct {
	Size    int
	Columns []Column
}

func (t *Tree) Format() {
	fmt.Printf("Object-Type: Tree  Size: %d\n", t.Size)
	for _, column := range t.Columns {
		fmt.Printf("%s %s %s\n", column.Type, column.Name, column.Hash)
	}
}

type Parent struct {
	Hash string
}

type Sign struct {
	Name      string
	Email     string
	TimeStamp time.Time
}

type Commit struct {
	Size      int
	Tree      string
	Parents   []Parent
	Author    Sign
	Committer Sign
	Message   string
}

func (c *Commit) Format() {
	fmt.Printf("Object-Type: Commit  Size: %d\n", c.Size)
	fmt.Println("Type     :", "Commit")
	fmt.Println("Size     :", c.Size)
	fmt.Println("Tree     :", c.Tree)
	for _, parent := range c.Parents {
		fmt.Println("Parents  :", parent)
	}
	fmt.Println("Author   :", c.Author)
	fmt.Println("Committer:", c.Committer)
	fmt.Println("Message  :", c.Message)
}

// func hash2Path(hash string) (string, error) {
// 	if len(hash) <= LENGTH_OF_HASH {
// 		return "", errors.New("Invalid Hash")
// 	}
// 	DirPath, ObjPath := hash[:2], hash[2:]
// 	hashPath := "/objects/" + DirPath + "/" + ObjPath
// 	return hashPath, nil
// }

// func hash2PathDir(hash string) (string, error) {
// 	if len(hash) <= LENGTH_OF_HASH {
// 		return "", errors.New("Invalid Hash")
// 	}
// 	DirPath := hash[:2]
// 	hashPath := "/objects/" + DirPath
// 	return hashPath, nil
// }

// func extract(zr io.Reader) (io.Reader, error) {
// 	return zlib.NewReader(zr)
// }

func GetGitHeader(buffer []byte) string {
	header := make([]byte, 1024)
	for _, buf := range buffer {
		if buf == 0 {
			break
		}
		header = append(header, buf)
	}
	return string(header)
}

func file2Buffer(f *os.File) []byte {
	buffer := make([]byte, 0)
	buf := make([]byte, 64)
	for {
		n, _ := (*f).Read(buf)
		if n == 0 {
			break
		}
		buffer = append(buffer, buf...)
	}
	return buffer
}

func Header3Content(buffer *[]byte) ([]byte, []byte, error) {
	datas := make([][]byte, 0)
	data := make([]byte, 0)
	for _, buf := range *buffer {
		if buf == 0 {
			if len(data) <= 1 {
				continue
			}
			datas = append(datas, data)
			data = make([]byte, 0)
		}
		data = append(data, buf)
	}
	datas = append(datas, data)
	if len(datas) < LENGTH_OF_DATA {
		return nil, nil, errors.New("fatal: Not a valid object name")
	}
	Header := datas[0]

	Content := make([]byte, 0)
	for _, data := range datas[1:] {
		Content = append(Content, data...)
	}
	return Header, Content, nil

}

func (c *Client) GetGitObject(hash string) ([]byte, error) {
	hash_rel_path := "/objects/" + hash[:2] + "/" + hash[2:]

	object_abs_path := c.Root + hash_rel_path
	f, err := os.Open(object_abs_path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buffer := make([]byte, 0)
	buf := make([]byte, 64)
	for {
		n, _ := (*f).Read(buf)
		if n == 0 {
			break
		}
		buffer = append(buffer, buf...)
	}

	extracting_buffer := bytes.NewBuffer(buffer)
	zlib_reader, err := zlib.NewReader(extracting_buffer)
	if err != nil {
		return nil, err
	}

	extracted_buffer, err := ioutil.ReadAll(zlib_reader)
	if err != nil {
		return nil, err
	}
	return extracted_buffer, nil
}

func Create3TObject(lineMeta string) string {
	hash := strings.Replace(lineMeta, "tree ", "", -1)
	hash = strings.ReplaceAll(hash, " ", "")
	return hash
}

func Create3PObject(lineMeta string) string {
	hash := strings.Replace(lineMeta, "parent ", "", -1)
	hash = strings.ReplaceAll(hash, " ", "")
	return hash
}

func Create3ACObject(lineMeta string) (Sign, error) {
	if ok := signRegexp.MatchString(lineMeta); !ok {
		return Sign{}, errors.New("NG")
	}
	sign1 := strings.SplitN(lineMeta, " <", 2)
	name := sign1[0]
	sign2 := strings.SplitN(sign1[1], "> ", 2)
	email := sign2[0]
	sign3 := strings.SplitN(sign2[1], " ", 2)
	unixTime, err := strconv.ParseInt(sign3[0], 10, 64)
	if err != nil {
		return Sign{}, err
	}
	var offsetHour, offsetMinute int
	if _, err := fmt.Sscanf(sign3[1], "+%02d%02d", &offsetHour, &offsetMinute); err != nil {
		return Sign{}, err
	}
	location := time.FixedZone(" ", 3600*offsetHour+60*offsetMinute)
	timestamp := time.Unix(unixTime, 0).In(location)
	time.Now().String()
	sign := Sign{
		Name:      name,
		Email:     email,
		TimeStamp: timestamp,
	}

	return sign, nil

}

func CreateCommitObject(buffer []byte) Commit {
	entries := make([][]byte, 0)
	entry := make([]byte, 0)

	for _, buf := range buffer {
		if buf == 0 {
			entries = append(entries, entry)
			entry = make([]byte, 0)
		}
		entry = append(entry, buf)
	}
	entries = append(entries, entry)

	for _, entry := range entries {
		fmt.Println(string(entry))
	}

	return Commit{}
}

func CreateCommitObject_v2(Header []byte, Content []byte) (Commit, error) {
	sizeStr := strings.Replace(string(Header), "commit ", "", -1)
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return Commit{}, err
	}
	cReader := strings.NewReader(string(Content[1:]))
	scanner := bufio.NewScanner(cReader)

	var commit Commit
	var parents []Parent
	commit.Size = size

	for scanner.Scan() {
		text := scanner.Text()
		columns := strings.SplitN(text, " ", 2)
		if len(columns) != NUM_OF_COMMIT_COLUMN {
			break
		}
		lineType := columns[0]
		lineMeta := columns[1]

		switch lineType {
		case "tree":
			hash := Create3TObject(lineMeta)
			commit.Tree = hash
		case "parent":
			var parent Parent
			hash := Create3PObject(lineMeta)
			parent.Hash = hash
			parents = append(parents, parent)
			commit.Parents = parents
		case "author":
			sign, err := Create3ACObject(lineMeta)
			if err != nil {
				continue
			}
			commit.Author = sign
		case "committer":
			sign, err := Create3ACObject(lineMeta)
			if err != nil {
				continue
			}
			commit.Committer = sign
		}
	}
	messages := make([]string, 0)
	for scanner.Scan() {
		messages = append(messages, scanner.Text())
	}
	message := strings.Join(messages, "\n")
	commit.Message = message

	return commit, nil

}

func CreateBlobObject(Header []byte, Content []byte) (Blob, error) {
	sizeStr := strings.Replace(string(Header), "blob ", "", -1)
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return Blob{}, err
	}
	blob := Blob{
		Size:    size,
		Content: Content,
	}
	return blob, nil
}

func CreateTreeObject(Header []byte, Content []byte) (Tree, error) {
	sizeStr := strings.Replace(string(Header), "tree ", "", -1)
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return Tree{}, err
	}
	buffers := make([][]byte, 0)
	buffer := make([]byte, 0)

	for _, cnt := range Content {
		if cnt == 0 {
			buffers = append(buffers, buffer)
			buffer = make([]byte, 0)
		}
		buffer = append(buffer, cnt)
	}
	buffers = append(buffers, buffer)

	lines := make([][]byte, 0)

	for _, buffer := range buffers {
		if len(buffer) <= 0 {
			continue
		}
		if len(buffer) >= 20 {
			hash := hex.EncodeToString(buffer[1:21])
			meta := buffer[21:]
			lines = append(lines, []byte(hash))
			lines = append(lines, meta)
		} else {
			meta := buffer[1:]
			lines = append(lines, meta)
		}
	}

	columns := make([]Column, 0)
	for n, line := range lines {
		if len(line) <= 0 {
			continue
		}
		if n%2 == 0 {
			if strings.HasPrefix(string(line), "40000") {
				name := strings.Replace(string(line), "40000 ", "", -1)
				column := Column{
					Type: "tree",
					Name: name,
					Hash: string(lines[n+1]),
				}
				columns = append(columns, column)
			} else {
				name := strings.Replace(string(line), "100644 ", "", -1)
				column := Column{
					Type: "blob",
					Name: name,
					Hash: string(lines[n+1]),
				}
				columns = append(columns, column)
			}
		}
	}
	tree := Tree{
		Size:    size,
		Columns: columns,
	}
	return tree, nil
}
