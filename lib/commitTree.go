package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func CommitTree(message string, hash string, client Client) (string, error) {
	ref_path, _ := GetHeadRef()

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

	if _, err := os.Stat(ref_path); err != nil {
		write_object, err := os.Create(ref_path)
		// write_object, err := os.Create(".bakibaki/refs/heads/master")
		// write_object, err := os.Create("/mnt/c/Users/81701/Documents/AtCoder/.bakibaki/refs/heads/master")
		if err != nil {
			fmt.Println("ERROR:1", err)
			return "", err
		}
		defer write_object.Close()

		if _, err = write_object.Write([]byte(hash)); err != nil {
			fmt.Println("ERROR:2", err)
			return "", err
		}
	} else {
		f, err := os.Open(ref_path)
		// f, err := os.Open("/mnt/c/Users/81701/Documents/AtCoder/.bakibaki/refs/heads/master")
		if err != nil {
			fmt.Println("ERROR:3", err)
			return "", err
		}
		defer f.Close()

		buffer, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Println("ERROR:4", err)
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

	return commit_hash, nil

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
}

type TatchedHead struct {
}

func GetHeadRef() (string, error) {
	current_dir, _ := os.Getwd()
	f, err := os.Open(current_dir + "/.bakibaki/HEAD")
	if err != nil {
		return "", err
	}
	defer f.Close()

	ref_buffer := make([]byte, 1024)
	if _, err := f.Read(ref_buffer); err != nil {
		return "", err
	}
	ref_string := string(ref_buffer)
	ref_string = strings.Replace(ref_string, "\n", "", -1)
	ref_string = strings.Replace(ref_string, "ref: ", "", 1)
	ref_string = strings.Replace(ref_string, ":", "", 1)
	ref_string = strings.Join(strings.Fields(ref_string), "")
	// return current_dir + "/.bakibaki/" + ref_string, nil
	ref_string = ".bakibaki/" + ref_string

	return ref_string, nil

}
