package lib

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
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

type GitObject interface {
	Format()
}

type Blob struct {
	Size    int
	Content []byte
}

func (b *Blob) Format() {
	fmt.Printf("Object-Type: Blob Size: %d\n", b.Size)
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
}

type Client struct {
	Root string
}

func hash2Path(hash string) (string, error) {
	if len(hash) <= LENGTH_OF_HASH {
		return "", errors.New("Invalid Hash")
	}
	DirPath, ObjPath := hash[:2], hash[2:]
	hashPath := "/objects/" + DirPath + "/" + ObjPath
	return hashPath, nil
}

func extract(zr io.Reader) (io.Reader, error) {
	return zlib.NewReader(zr)
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
		return nil, nil, errors.New("Header or Contents is empty")
	}
	Header := datas[0]

	Content := make([]byte, 0)
	for _, data := range datas[1:] {
		Content = append(Content, data...)
	}
	return Header, Content, nil

}

func (c *Client) GetGitObject(hash string) ([]byte, error) {
	hashPath, err := hash2Path(hash)
	if err != nil {
		return nil, err
	}
	ObjectPath := c.Root + hashPath
	f, err := os.Open(ObjectPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buffer := file2Buffer(f)

	Ebuffer := bytes.NewBuffer(buffer)

	Zf, err := extract(Ebuffer)
	if err != nil {
		return nil, err
	}

	Zbuffer, err := ioutil.ReadAll(Zf)
	if err != nil {
		return nil, err
	}
	return Zbuffer, nil
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

func CreateCommitObject(Header string, Content string) (Commit, error) {
	sizeStr := strings.Replace(Header, "commit ", "", -1)
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


