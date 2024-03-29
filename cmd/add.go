/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/aoimaru/bakibaki/lib"
	"github.com/spf13/cobra"
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
		BakiBakiRootPath, err := lib.FindBakiBakiRoot(current)
		if err != nil {
			fmt.Println(err)
		}

		client := lib.Client{
			Root: BakiBakiRootPath,
		}
		path := args[0]
		buffer, hash, err := client.CreateBlobFile(path)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(buffer), hash)

		// indexファイルをオブジェクトとして取得
		index_path := client.GetIndexPath()
		index, err := client.GetIndexObject(index_path)
		if err != nil {
			fmt.Println(err)
		}

		new_index := index.UpdateIndex(path, hash)
		index_buffer := new_index.AsByte()
		if err := index_buffer.ToFile(client); err != nil {
			fmt.Println(err)
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
