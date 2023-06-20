/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/aoimaru/bakibaki/lib"
	"github.com/spf13/cobra"
)

const (
	NUM_OF_HO_ARGS = 1
)

// hashObjectCmd represents the hashObject command
var hashObjectCmd = &cobra.Command{
	Use:   "hashObject",
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
		path := args[0]
		buffer, hash, err := client.CreateBlobFile(path)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(buffer), hash)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		/** 引数のバリデーションを行うことができる */
		if len(args) < NUM_OF_HO_ARGS {
			return errors.New("requires args")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(hashObjectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hashObjectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hashObjectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
