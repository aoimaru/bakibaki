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

// lsFilesCmd represents the lsFiles command
var lsFilesCmd = &cobra.Command{
	Use:   "lsFiles",
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
		client := lib.Client{
			Root: GitRootPath,
		}
		indexPath := client.GetIndexPath()
		fmt.Println(indexPath)
		// index, err := lib.GetIndexObject(indexPath)
		index, err := lib.GetIndexObject("/mnt/c/Users/81701/Documents/AtCoder/subsubIndex")
		if err != nil {
			fmt.Println(err)
		}
		for _, entry := range (*index).Entries {
			fmt.Println(entry.Name, entry.Hash)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsFilesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsFilesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsFilesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
