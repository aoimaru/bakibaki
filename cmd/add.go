/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"errors"
	"path/filepath"


	"github.com/spf13/cobra"

	"github.com/aoimaru/bakibaki/lib"
)

const (
	NUM_OF_ADD_ARGS = 1
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {



		current, _ := os.Getwd()
		GitRootPath, err := lib.FindBakiBakiRoot(current)
		if err != nil {
			fmt.Println(err)
		}
		
		client := lib.Client{
			Root: GitRootPath,
		}

		RepRootPath := GitRootPath

		/** この部分はSATD*/
		if strings.HasSuffix(GitRootPath, "/.git") {
			RepRootPath = strings.Replace(GitRootPath, "/.git", "", -1)+"/"
		}
		if strings.HasSuffix(GitRootPath, "/.bakibaki") {
			RepRootPath = strings.Replace(GitRootPath, "/.bakibaki", "", -1)+"/"
		}
		fmt.Println(RepRootPath)

		indexPath := client.GetIndexPath()
		index, err := lib.GetIndexObject(indexPath)
		if err != nil {
			fmt.Println(err)
		}


		for _, path := range args {
			filePath, err := filepath.Abs(path)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if strings.Contains(filePath, "/.git/") {
				continue
			}
			if strings.Contains(filePath, "/.bakibaki/") {
				continue
			}

			_, hash, err := client.CreateBlobFile(filePath)
			if err != nil {
				fmt.Println(err)
				continue
			}

			name := strings.Replace(filePath, RepRootPath, "", -1)

			Nindex, _, err := lib.UpdateIndex(index, name, hash, &client)
			if err != nil {
				fmt.Println(err)
			}
			err = lib.WriteIndex(Nindex, indexPath)
			if err != nil {
				fmt.Println(err)
			}

		}


	},
	Args: func(cmd *cobra.Command, args []string) error {
		/** 引数のバリデーションを行うことができる */
		if len(args) < NUM_OF_ADD_ARGS {
			return errors.New("requires args")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
