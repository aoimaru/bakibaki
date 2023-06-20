/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/aoimaru/bakibaki/lib"
	"github.com/spf13/cobra"
)

// writeTreeCmd represents the writeTree command
var writeTreeCmd = &cobra.Command{
	Use:   "writeTree",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// BakiBakiリポジトリのルートパスを取得
		current, _ := os.Getwd()
		// GitRootPath, err := lib.FindBakiBakiRoot(current)
		GitRootPath, err := lib.FindGitRoot(current)
		if err != nil {
			fmt.Println(err)
		}
		client := lib.Client{
			Root: GitRootPath,
		}

		// indexファイルをオブジェクトとして取得
		indexPath := client.GetIndexPath()
		fmt.Println(indexPath)
		// index, err := lib.GetIndexObject(indexPath)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// // util.CreateTree(index)
		// fmt.Println(index)
	},
}

func init() {
	rootCmd.AddCommand(writeTreeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// writeTreeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// writeTreeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
