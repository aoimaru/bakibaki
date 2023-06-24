/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/aoimaru/bakibaki/lib"
	"github.com/aoimaru/bakibaki/util"
	"github.com/spf13/cobra"
	// "github.com/aoimaru/bakibaki/test"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		current_dir, _ := os.Getwd()
		BakiBakiRootPath, err := lib.FindBakiBakiRoot(current_dir)
		if err != nil {
			fmt.Println(err)
		}

		client := lib.Client{
			Root: BakiBakiRootPath,
		}
		fmt.Println(client)
		root_dir := current_dir + "/" + args[0]
		if args[0] == "." {
			root_dir = current_dir
		}

		working_file_paths, err := util.WalkingDir(root_dir)
		if err != nil {
			fmt.Println(err)
		}
		// indexファイルをオブジェクトとして取得
		index_path := client.GetIndexPath()
		index, err := client.GetIndexObject(index_path)
		if err != nil {
			fmt.Println(err)
		}

		for _, working_file_path := range working_file_paths {
			fmt.Println("working:", working_file_path)
			buffer, hash, err := client.CreateBlobFile(working_file_path)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(buffer), hash)
			index = index.UpdateIndex(working_file_path, hash)
		}
		index_buffer := index.AsByte()
		if err := index_buffer.ToFile(client); err != nil {
			fmt.Println(err)
		}

	},
	Args: func(cmd *cobra.Command, args []string) error {
		/** 引数のバリデーションを行うことができる */
		if len(args) < 1 {
			return errors.New("requires args")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
