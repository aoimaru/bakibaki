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

// writeTreeTestCmd represents the writeTreeTest command
var writeTreeTestCmd = &cobra.Command{
	Use:   "writeTreeTest",
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
		index_path := client.GetIndexPath()
		fmt.Println(index_path)

		index, err := client.GetIndexObject(index_path)
		if err != nil {
			fmt.Println(err)
		}
		nodes := index.CreateNodes()
		for _, node := range nodes {
			fmt.Printf("node: %+v\n", node)
		}
	},
}

func init() {
	rootCmd.AddCommand(writeTreeTestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// writeTreeTestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// writeTreeTestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
