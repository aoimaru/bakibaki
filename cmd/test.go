/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/aoimaru/bakibaki/lib"
	"github.com/aoimaru/bakibaki/test"
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
		current, _ := os.Getwd()
		GitRootPath, _ := lib.FindGitRoot(current)
		client := lib.Client{
			Root: GitRootPath,
		}
		indexPath := client.GetIndexPath()
		index, _ := lib.GetIndexObject(indexPath)
		
		// newEntry, filePath, _ := lib.UpdateIndex(index, "bakibaki.py", "0e50249a75625c1b02a04103cca4a3027128da4c", &client)
		// fmt.Println("filePath->", filePath)
		// err := lib.WriteIndex(newEntry, "./index1")
		// if err != nil {
		// 	fmt.Println(err)
		// }

		// lib.TestIndex()

		
		fmt.Println("---> test!!")
		test.WriteIndexHeaderTest(index)
		
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
