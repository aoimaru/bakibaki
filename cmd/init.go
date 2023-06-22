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

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")
		current_dir, _ := os.Getwd()

		if _, err := os.Stat(current_dir + "/.bakibaki"); err == nil {
			fmt.Println(current_dir + "/.bakibaki")
			_ = os.RemoveAll(current_dir + "/.bakibaki")
		}

		if err := os.MkdirAll(current_dir+"/.bakibaki/objects", 1755); err != nil {
			fmt.Println(err)
		}
		if err := os.MkdirAll(current_dir+"/.bakibaki/refs/heads", 1755); err != nil {
			fmt.Println(err)
		}

		// BakiBakiクライアントを作成
		GitRootPath, err := lib.FindBakiBakiRoot(current_dir)
		if err != nil {
			fmt.Println(err)
		}
		client := lib.Client{
			Root: GitRootPath,
		}

		// indexファイルのファイルパスを取得
		// index_path := client.GetIndexPath()

		// indexファイルが存在しない場合は, エントリーが空のindexファイルを作成する
		index := lib.InitIndexObject()
		index_buffer := index.AsByte()
		index_buffer.ToFile(client)

		if err = lib.CreateHEAD(); err != nil {
			fmt.Println(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
