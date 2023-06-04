/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/aoimaru/bakibaki/lib"
	"github.com/aoimaru/bakibaki/util"
	"github.com/spf13/cobra"
)

const (
	NUM_OF_ADD_ARGS = 1
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		current, _ := os.Getwd()
		BakiBakiRootPath, err := lib.FindBakiBakiRoot(current)
		GitRootPath, err := lib.FindGitRoot(current)
		if err != nil {
			fmt.Println(err)
		}

		bakibaki_client := lib.Client{
			Root: BakiBakiRootPath,
		}
		name := args[0]
		buffers, hash, err := bakibaki_client.CreateBlobFile(name)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(buffers), hash)

		// indexファイルをオブジェクトとして取得

		client := lib.Client{
			Root: GitRootPath,
		}

		indexPath := client.GetIndexPath()
		index, err := lib.GetIndexObject(indexPath)
		if err != nil {
			fmt.Println(err)
		}

		// name := "typical90/032_TLE-2.py"
		// hash := "f131105cf1b5940a07b76f2608ea605f1ebcf2c7"

		filePath := current + "/" + name
		// fmt.Println(hash)

		var sysC syscall.Stat_t
		syscall.Stat(filePath, &sysC)

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			fmt.Println(filePath, err)
			return
		}

		// fmt.Println(fileInfo)

		oct := fmt.Sprintf("%o", uint32(sysC.Mode))
		num, err := strconv.ParseUint(oct, 10, 32)
		if err != nil {
			fmt.Println(err)
		}
		mode := uint32(num)

		new_entry := lib.Entry{
			CTime: fileInfo.ModTime(),
			MTime: fileInfo.ModTime(),
			Dev:   uint32(sysC.Dev),
			Inode: uint32(sysC.Ino),
			Mode:  mode,
			Uid:   sysC.Uid,
			Gid:   sysC.Gid,
			Size:  uint32(sysC.Size),
			Hash:  hash,
			Name:  name,
		}

		var new_index lib.Index

		for _, entry := range index.Entries {
			if entry.Name == name {
				continue
			}
			if entry.Hash == hash {
				continue
			}
			fmt.Printf("%+v\n", entry)
			new_index.Entries = append(new_index.Entries, entry)
		}
		new_index.Entries = append(new_index.Entries, new_entry)

		// fmt.Println(len(index.Entries))
		for _, entry := range new_index.Entries {
			fmt.Println("Index", entry)
		}
		// fmt.Println(len(new_index.Entries))

		new_index.Dirc = "DIRC"
		new_index.Version = 2
		new_index.Number = uint32(len(new_index.Entries))

		buffer := make([]byte, 0)

		dirc := []byte(new_index.Dirc)
		version := util.Element2byte32(new_index.Version)
		number := util.Element2byte32(new_index.Number)

		buffer = append(buffer, dirc...)
		buffer = append(buffer, version...)
		buffer = append(buffer, number...)

		for _, entry := range new_index.Entries {

			c_unix := entry.CTime.Unix()
			buf_c_unix := util.Element2byte32(uint32(c_unix))
			buffer = append(buffer, buf_c_unix...)
			buffer = append(buffer, buf_c_unix...)

			m_unix := entry.MTime.Unix()
			buf_m_unix := util.Element2byte32(uint32(m_unix))
			buffer = append(buffer, buf_m_unix...)
			buffer = append(buffer, buf_m_unix...)

			dev := entry.Dev
			buf_dev := util.Element2byte32(uint32(dev))
			buffer = append(buffer, buf_dev...)

			inode := entry.Inode
			buffer_inode := util.Element2byte32(uint32(inode))
			buffer = append(buffer, buffer_inode...)

			mode := entry.Mode
			buffer_mode := util.Element2byte32(uint32(mode))
			buffer = append(buffer, buffer_mode...)

			uid := entry.Uid
			buffer_uid := util.Element2byte32(uint32(uid))
			buffer = append(buffer, buffer_uid...)

			gid := entry.Gid
			buffer_gid := util.Element2byte32(uint32(gid))
			buffer = append(buffer, buffer_gid...)

			size := entry.Size
			buffer_size := util.Element2byte32(uint32(size))
			buffer = append(buffer, buffer_size...)

			bHash, err := hex.DecodeString(entry.Hash)
			if err != nil {
				fmt.Println("ココ？")
				continue
			}
			buffer = append(buffer, bHash...)

			bnSize := make([]byte, 2)
			binary.BigEndian.PutUint16(bnSize, uint16(len(entry.Name)))
			buffer = append(buffer, bnSize...)

			bName := []byte(entry.Name)
			buffer = append(buffer, bName...)

			var sw uint64
			sw = 62

			padding := util.GetPaddingSize(sw + uint64(len(bName)))
			bPadding := make([]byte, padding)
			buffer = append(buffer, bPadding...)
			fmt.Println(entry)
			fmt.Println("OK")
		}

		w, err := os.Create(indexPath)
		// w, err := os.Create("/mnt/c/Users/81701/Documents/AtCoder/subsubsubIndex")
		if err != nil {
			fmt.Println(err)
		}
		defer w.Close()

		count, err := w.Write(buffer)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("write %d bytes\n", count)
		// fmt.Println("TEST")

	},
	Args: func(cmd *cobra.Command, args []string) error {
		/** 引数のバリデーションを行うことができる */
		if len(args) < NUM_OF_ADD_ARGS {
			return errors.New("requires args")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
