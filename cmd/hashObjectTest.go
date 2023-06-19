/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/aoimaru/bakibaki/lib"
	"github.com/aoimaru/bakibaki/util"
	"github.com/spf13/cobra"
)

// hashObjectTestCmd represents the hashObjectTest command
var hashObjectTestCmd = &cobra.Command{
	Use:   "hashObjectTest",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		current, _ := os.Getwd()
		// GitRootPath, err := lib.FindGitRoot(current)
		GitRootPath, err := lib.FindBakiBakiRoot(current)
		if err != nil {
			fmt.Println(err)
		}

		hash := args[0]
		client := lib.Client{
			Root: GitRootPath,
		}
		hash_object, err := client.GetGitObject(hash)
		if err != nil {
			fmt.Println(err)
		}
		header, _ := util.GetGitObjectHeader(&hash_object)
		fmt.Println(header)
	},
}

func init() {
	rootCmd.AddCommand(hashObjectTestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hashObjectTestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hashObjectTestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
