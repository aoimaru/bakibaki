/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/aoimaru/bakibaki/lib"
	// "github.com/aoimaru/bakibaki/test"

	"syscall"
	"strconv"
	"time"
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
		// _, _ = lib.GetIndexObject(indexPath)

		for _, entry := range (*index).Entries {
			fmt.Println(entry.Name, entry.Hash, entry.Size)
		}
		
		fmt.Println("")
		// test.WriteIndexHeaderTest(index)

		name := "typical90/032_TLE-2.py"
		hash := "f131105cf1b5940a07b76f2608ea605f1ebcf2c7"

		filePath := current+"/"+name
		fmt.Println(hash)

		var sysC syscall.Stat_t
		syscall.Stat(filePath, &sysC)

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			fmt.Println(filePath, err)
		}

		// fmt.Println(fileInfo)

		oct := fmt.Sprintf("%o", uint32(sysC.Mode))
		num, err := strconv.ParseUint(oct, 10, 32)
		if err != nil {
			fmt.Println(err)
		}
		mode := uint32(num)

		type Entry struct {
			cTime time.Time
			mTime time.Time
			Dev   uint32
			Inode uint32
			Mode  uint32
			Uid   uint32
			Gid   uint32
			Size  uint32
			Hash  string
			Name  string
		}
		
		Nentry := Entry{
			cTime: fileInfo.ModTime(),
			mTime: fileInfo.ModTime(),
			Dev:   uint32(sysC.Dev),
			Inode: uint32(sysC.Ino),
			Mode:  mode,
			Uid:   sysC.Uid,
			Gid:   sysC.Gid,
			Size:  uint32(sysC.Size),
			Hash:  hash,
			Name:  name,
		}

		type Index2 struct {
			Dirc    string
			Version uint32
			Number  uint32
			Entries []Entry
		}

		fmt.Printf("%+v", Nentry)
		var index2 Index2
		index2.Dirc = "DIRC"
		index2.Version = 2
		index2.Number = 1

		index2.Entries = append(index2.Entries, Nentry)



		

		
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
