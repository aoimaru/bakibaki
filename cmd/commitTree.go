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

// commitTreeCmd represents the commitTree command
var commitTreeCmd = &cobra.Command{
	Use:   "commitTree",
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
		BakiBakiRootPath, err := lib.FindBakiBakiRoot(current)
		if err != nil {
			fmt.Println(err)
		}

		client := lib.Client{
			Root: BakiBakiRootPath,
		}
		commit_hash, err := lib.CommitTree(message, hash, client)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("commit hash:", commit_hash)
	},
}

func init() {
	rootCmd.AddCommand(commitTreeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commitTreeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commitTreeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	commitTreeCmd.Flags().StringP("message", "m", "", "set file message")
	commitTreeCmd.Flags().StringP("hash", "s", "", "set file hash")
}
