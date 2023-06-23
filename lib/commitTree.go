package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

func CommitTree(message string, hash string, client Client) (string, error) {
	var commit Commit
	commit.Size = 119
	commit.Tree = hash

	var sign Sign
	sign.Name = "aoi nakamura"
	sign.Email = "hello@world.com"
	sign.TimeStamp = time.Now()

	commit.Author = sign
	commit.Committer = sign

	commit.Message = message

	ref_head, err := GetHeadRef()
	if err != nil {
		return "", nil
	}

	switch v := ref_head.(type) {
	case DetachedHead:
		return v.Head, nil
	case TatchedHead:
		if _, err := os.Stat(v.Head); err != nil {
			err = ioutil.WriteFile(v.Head, []byte(hash), 0664)
			if err != nil {
				return "", err
			}
			fmt.Println(v.Head)
		} else {
			f, err := os.Open(v.Head)
			if err != nil {
				return "", err
			}
			defer f.Close()
			buffer, err := ioutil.ReadAll(f)
			if err != nil {
				fmt.Println(err)
				return "", err
			}
			parent_hash := string(buffer)
			commit.Parents = append(commit.Parents, Parent{Hash: parent_hash})
		}
		commit.Format()
		commit_buffer := commit.AsByte()
		commit_hash, err := commit_buffer.ToFile(client)
		if err != nil {
			return "", err
		}
		fmt.Println(commit_hash)
		return commit_hash, nil
	default:
		return "", err
	}
}

func (c *Commit) AsByte() CommitBuffer {

	buffer := make([]byte, 0)
	buffer = append(buffer, []byte("commit 199")...)
	buffer = append(buffer, 0)

	if len(c.Parents) > 0 {
		for _, parent := range c.Parents {
			parent_string := "parent " + parent.Hash
			buffer = append(buffer, []byte(parent_string)...)
			buffer = append(buffer, 0)
		}
	}

	tree_string := "tree " + c.Tree
	buffer = append(buffer, []byte(tree_string)...)
	buffer = append(buffer, 0)

	author_string := "author " + c.Author.Name + " " + c.Author.Email + " " + c.Author.TimeStamp.String()
	buffer = append(buffer, []byte(author_string)...)
	buffer = append(buffer, 0)

	committer_string := "committer " + c.Committer.Name + " " + c.Committer.Email + " " + c.Committer.TimeStamp.String()
	buffer = append(buffer, []byte(committer_string)...)
	buffer = append(buffer, 0)

	buffer = append(buffer, []byte(c.Message)...)

	return CommitBuffer{Buffer: buffer}
}

type Head interface {
}

type DetachedHead struct {
	Head string
}

type TatchedHead struct {
	Head string
}

func GetHeadRef() (Head, error) {
	current_dir, _ := os.Getwd()
	f, err := os.Open(current_dir + "/.bakibaki/HEAD")
	if err != nil {
		return DetachedHead{}, err
	}
	defer f.Close()
	ref_buffer := make([]byte, 1024)
	if _, err := f.Read(ref_buffer); err != nil {
		return "", err
	}
	ref_string := string(ref_buffer)

	re := regexp.MustCompile(`ref: refs/heads/(\w+)`)

	if re.MatchString(ref_string) {
		ref := make([]byte, 0)
		for _, ref_buf := range ref_buffer {
			if ref_buf == 0 {
				break
			}
			ref = append(ref, ref_buf)
		}
		ref_string = string(ref)
		ref_string = strings.Replace(ref_string, "\n", "", -1)
		ref_string = strings.Replace(ref_string, "ref: ", "", 1)
		ref_string = strings.Replace(ref_string, ":", "", 1)

		current_dir, _ := os.Getwd()
		ref_string = current_dir + "/.bakibaki/" + ref_string
		ref_string = strings.Join(strings.Fields(ref_string), "")
		return TatchedHead{Head: ref_string}, nil
	} else {
		return DetachedHead{Head: ref_string}, nil
	}

}
