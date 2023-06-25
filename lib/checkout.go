package lib

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/aoimaru/bakibaki/util"
)

func (c *Client) GetHeadHash() (string, error) {
	head_path := c.Root + "/HEAD"
	f, err := os.Open(head_path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	buffer, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	head_string := string(buffer)
	re := regexp.MustCompile(`ref: refs/heads/(\w+)`)
	if re.MatchString(head_string) {
		head_string = strings.Replace(head_string, "\n", "", -1)
		head_string = strings.Replace(head_string, "ref: ", "", 1)
		head_string = strings.Replace(head_string, ":", "", 1)
		fmt.Println(c.Root + "/" + head_string)
		ref_f, err := os.Open(c.Root + "/" + head_string)
		if err != nil {
			fmt.Println("ERROR:1")
			return "", err
		}
		defer ref_f.Close()
		ref_hash, err := ioutil.ReadAll(ref_f)
		if err != nil {
			fmt.Println("ERROR:2")
			return "", err
		}
		return string(ref_hash), nil
	} else {
		return head_string, nil
	}
}

func (c *Client) CreateBranch(branch_name string, hash string) error {
	branches, err := util.GetBrabches(c.Root + "/refs/heads")
	flag := false
	for _, branch := range branches {
		if branch == branch_name {
			flag = true
			break
		}
	}
	if flag {
		return errors.New("branch is exist")
	}
	branch_path := c.Root + "/refs/heads/" + branch_name
	write_object, err := os.Create(branch_path)
	if err != nil {
		fmt.Println("ERROR:1", err)
		return err
	}
	defer write_object.Close()

	if _, err = write_object.Write([]byte(hash)); err != nil {
		fmt.Println("ERROR:2", err)
		return err
	}

	head_path := c.Root + "/HEAD"
	f, err := os.Create(head_path)
	if err != nil {
		return err
	}
	defer f.Close()
	ref_name := "ref: refs/heads/" + branch_name
	if _, err = f.Write([]byte(ref_name)); err != nil {
		fmt.Println("ERROR:2", err)
		return err
	}

	return nil
}

type Checkout struct {
	Path string
	Hash string
}

func (c *Client) GetCommitHashFromName(branch_name string) (string, error) {
	branch_path := c.Root + "/refs/heads/" + branch_name
	buffer, err := os.ReadFile(branch_path)
	if err != nil {
		return "", err
	}
	return string(buffer), nil
}

func (c *Client) CreateCheckoutObject(commit string, file_path string) (Checkout, error) {
	if _, err := os.Stat(file_path); err != nil {
		return Checkout{}, err
	}
	root_path := c.Root + "/refs/heads"
	branches, _ := util.GetBrabches(root_path)
	for _, branch := range branches {
		fmt.Println("branch:", branch)
		commit_id, _ := c.GetCommitHashFromName(branch)
		if commit == branch {
			return Checkout{Hash: commit_id, Path: file_path}, nil
		}
	}
	commit_ids := c.GetCommitObjectID()
	for _, commit_id := range commit_ids {
		if commit == commit_id {
			return Checkout{Hash: commit, Path: file_path}, nil
		}
	}

	return Checkout{}, nil
}

func (co *Checkout) RollBackIndex(client *Client) {
	buffer, err := client.GetGitObject(co.Hash)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(buffer))
	tree_hash, _ := client.GetCommitTreeHash(buffer)
	fmt.Println(tree_hash)
	_ = client.GetTreeObject(tree_hash)

}
