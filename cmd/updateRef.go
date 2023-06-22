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

// updateRefCmd represents the updateRef command
var updateRefCmd = &cobra.Command{
	Use:   "updateRef",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		head, err := cmd.Flags().GetString("head")
		if err != nil {
			fmt.Println("set file head")
			return
		}
		if head == "" {
			fmt.Println("set file head")
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

		err = client.UpdateRef(head, hash)
		if err != nil {
			fmt.Println(err)
		}

		// current_dir, _ := os.Getwd()
		// head_path := current_dir + "/.bakibaki/" + head
		// fmt.Println(head_path)
	},
}

func init() {
	rootCmd.AddCommand(updateRefCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateRefCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateRefCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	updateRefCmd.Flags().StringP("head", "n", "", "set file head")
	updateRefCmd.Flags().StringP("hash", "s", "", "set file hash")
}
