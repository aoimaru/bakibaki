/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"

	"errors"
	"github.com/aoimaru/bakibaki/lib"
	"strings"
)

const (
	NUM_OF_ARGS = 1
)

// catFileCmd represents the catFile command
var catFileCmd = &cobra.Command{
	Use:   "catFile",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		current, _ := os.Getwd()
		GitRootPath, err := lib.FindGitRoot(current)
		// GitRootPath, err := lib.FindBakiBakiRoot(current)
		if err != nil {
			fmt.Println(err)
		}
		
		hash := args[0]
		client := lib.Client{
			Root: GitRootPath,
		}
		GObuffer, err := client.GetGitObject(hash)
		if err != nil {
			fmt.Println("ERROR:1")
			fmt.Println(err)
		}
		Header, Content, err := lib.Header3Content(&GObuffer)
		if err != nil {
			fmt.Println("ERROR:2")
			fmt.Println(err)
		}

		if strings.HasPrefix(string(Header), "blob ") {
			blob, err := lib.CreateBlobObject(Header, Content)
			if err != nil {
				fmt.Println(err)
			}
			blob.Format()
		} else if strings.HasPrefix(string(Header), "tree ") {
			tree, err := lib.CreateTreeObject(Header, Content)
			if err != nil {
				fmt.Println(err)
			}
			tree.Format()
		} else if strings.HasPrefix(string(Header), "commit ") {
			commit, err := lib.CreateCommitObject(Header, Content)
			if err != nil {
				fmt.Println(err)
			
			}
			commit.Format()
		}
 	},
	Args: func(cmd *cobra.Command, args []string) error {
		/** 引数のバリデーションを行うことができる */
		if len(args) < NUM_OF_ARGS {
			return errors.New("requires args")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(catFileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// catFileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// catFileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
