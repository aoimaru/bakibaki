/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	// "strings"
	"errors"
	"github.com/spf13/cobra"

	"github.com/aoimaru/bakibaki/lib"
)
const (
	NUM_OF_UI_ARGS = 2
)


// updateIndexCmd represents the updateIndex command
var updateIndexCmd = &cobra.Command{
	Use:   "updateIndex",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Println("set file name")
			return
		}
		if name == "" {
			fmt.Println("set file name")
			return
		}

		hash, err := cmd.Flags().GetString("hash")
		if err != nil {
			fmt.Println("set file hash")
			return
		}
		if hash == "" {
			fmt.Println("set file hash")
			return
		}

		current, _ := os.Getwd()
		GitRootPath, err := lib.FindGitRoot(current)
		if err != nil {
			fmt.Println(err)
		}
		client := lib.Client{
			Root: GitRootPath,
		}
		indexPath := client.GetIndexPath()
		index, err := lib.GetIndexObject(indexPath)
		if err != nil {
			fmt.Println(err)
		}


		
		lib.UpdateIndex(index, "Hello", "World")

	},
	Args: func(cmd *cobra.Command, args []string) error {
		/** 引数のバリデーションを行うことができる */
		if len(args) < NUM_OF_UI_ARGS {
			return errors.New("requires args")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateIndexCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateIndexCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	updateIndexCmd.Flags().StringP("name", "n", "", "set file name")
	updateIndexCmd.Flags().StringP("hash", "s", "", "set file hash")
}
