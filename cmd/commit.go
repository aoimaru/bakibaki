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

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		message, err := cmd.Flags().GetString("message")
		if err != nil {
			fmt.Println("set file message")
			return
		}
		if message == "" {
			fmt.Println("set file message")
			return
		}
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
			if (*node).Path == "root" {
				root_tree_hash := lib.WriteTree(node, &index)
				fmt.Println("root tree hash:", root_tree_hash)

				commit_hash, err := lib.CommitTree(message, root_tree_hash, client)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("commit hash:", commit_hash)
				// head := "refs/heads/master"
				head, err := lib.GetHeadRef()
				if err != nil {
					fmt.Println(err)
				}
				err = client.UpdateRef(head, commit_hash)
				if err != nil {
					fmt.Println(err)
				}
				break
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	commitCmd.Flags().StringP("message", "m", "", "set file message")
}
