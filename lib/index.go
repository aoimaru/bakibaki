package lib

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
	// "reflect"
	"encoding/hex"
	// "io/ioutil"
)

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

type Index struct {
	Dirc    string
	Version uint32
	Number  uint32
	Entries []Entry
}

// Bytes2uint converts []byte to uint64
func Bytes2uint(bytes []byte) uint64 {
	padding := make([]byte, 8-len(bytes))
	i := binary.BigEndian.Uint64(append(padding, bytes...))
	return i
}

func Bytes2Uint32(bytes []byte) uint32 {
	padding := make([]byte, 4-len(bytes))
	i := binary.BigEndian.Uint32(append(padding, bytes...))
	return i
}

func Bytes2Uint16(bytes []byte) uint16 {
	padding := make([]byte, 2-len(bytes))
	i := binary.BigEndian.Uint16(append(padding, bytes...))
	return i
}

// Bytes2str converts []byte to string("00 00 00 00 00 00 00 00")
func Bytes2str(bytes []byte) string {
	strs := []string{}
	for _, b := range bytes {
		strs = append(strs, fmt.Sprintf("%02x", b))
	}
	return strings.Join(strs, " ")
}

func ChatGPT(bytes []byte) (uint32, error) {
	dec := binary.BigEndian.Uint32(bytes)
	oct := fmt.Sprintf("%o", dec)
	num, err := strconv.ParseUint(oct, 10, 32)
	if err != nil {
		return 0, err
	}
	uint32Num := uint32(num)
	return uint32Num, nil
}

func GetPaddingSize(had uint64) uint64 {
	Rem := had % 8
	return 8 - Rem
}

func GetUnixTime(sTime uint64) (time.Time, error) {
	unixTime := int64(sTime)
	var offsetHour, offsetMinute int
	if _, err := fmt.Sscanf("+0900", "+%02d%02d", &offsetHour, &offsetMinute); err != nil {
		return time.Time{}, err
	}
	location := time.FixedZone(" ", 3600*offsetHour+60*offsetMinute)
	timestamp := time.Unix(unixTime, 0).In(location)
	time.Now().String()
	return timestamp, nil
}

func ten2eight(mode uint64) uint64 {
	tmode := int64(mode)
	emode := strconv.FormatInt(tmode, 4)
	rmode, _ := strconv.ParseInt(emode, 10, 32)
	return uint64(rmode)
}

func GetIndexFile(file_path string) ([]byte, error) {
	f, err := os.Open(file_path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buffer, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func CreateIndex(buffer []byte) (*Index, error) {
	dirc := string(buffer[0:4])
	if dirc != "DIRC" {
		return nil, errors.New("NOT INDEX FILE")
	}

	version := Bytes2Uint32(buffer[4:8])
	if version != 2 {
		err := errors.New("Invalid Version Error")
		return nil, err
	}

	enum := Bytes2Uint32(buffer[8:12])

	buffer = buffer[12:]

	var index Index

	index.Dirc = dirc
	index.Version = version
	index.Number = enum

	var count uint32
	count = 0
	for {
		if count >= enum {
			break
		}
		count++

		ctime, err := GetUnixTime(Bytes2uint(buffer[0:4]))
		if err != nil {
			fmt.Println(err)
			continue
		}
		_ = Bytes2str(buffer[4:8])

		mtime, err := GetUnixTime(Bytes2uint(buffer[8:12]))
		if err != nil {
			fmt.Println(err)
			continue
		}
		_ = Bytes2str(buffer[12:16])

		dev := Bytes2Uint32(buffer[16:20])
		inode := Bytes2Uint32(buffer[20:24])
		mode, err := ChatGPT(buffer[24:28])
		if err != nil {
			continue
		}
		uid := Bytes2Uint32(buffer[28:32])
		gid := Bytes2Uint32(buffer[32:36])
		size := Bytes2Uint32(buffer[36:40])
		hash := hex.EncodeToString(buffer[40:60])
		nsize := Bytes2uint(buffer[60:62])
		name := string(buffer[62 : 62+nsize])

		entry := Entry{
			cTime: ctime,
			mTime: mtime,
			Dev:   dev,
			Inode: inode,
			Mode:  mode,
			Uid:   uid,
			Gid:   gid,
			Size:  size,
			Hash:  hash,
			Name:  name,
		}

		index.Entries = append(index.Entries, entry)

		padding := GetPaddingSize(62 + nsize)
		offset := 62 + nsize + padding
		buffer = buffer[offset:]
	}

	return &index, nil
}

func GetIndexObject(file_path string) (*Index, error) {
	buffer, err := GetIndexFile(file_path)
	if err != nil {
		return nil, err
	}
	index, err := CreateIndex(buffer)
	if err != nil {
		return nil, err
	}
	return index, nil
}

func UpdateIndex(index *Index, name string, hash string, client *Client) (*Index, string, error) {
	current, _ := os.Getwd()
	filePath := current + "/" + name

	var sysC syscall.Stat_t
	syscall.Stat(filePath, &sysC)

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, "", err
	}
	oct := fmt.Sprintf("%o", uint32(sysC.Mode))
	num, err := strconv.ParseUint(oct, 10, 32)
	if err != nil {
		return nil, "", err
	}
	mode := uint32(num)
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

	fmt.Printf("%+v\n", Nentry)

	var Nindex Index

	Nindex.Dirc = (*index).Dirc
	Nindex.Version = (*index).Version
	Nindex.Number = (*index).Number

	for _, entry := range (*index).Entries {
		if entry.Name == name {
			continue
		}
		if entry.Hash == hash {
			continue
		}
		Nindex.Entries = append(Nindex.Entries, entry)
		Nentry.cTime = entry.cTime /** ここでファイルの作成時間を遺伝*/
	}

	Nindex.Entries = append(Nindex.Entries, Nentry)

	return &Nindex, filePath, nil
}

func Uint642Byte(ui uint32) []byte {
	bt := make([]byte, 4)
	binary.BigEndian.PutUint32(bt, ui)
	return bt
}

func WriteIndex(index *Index, file_path string) error {
	buffer := make([]byte, 0)

	bDirc := []byte((*index).Dirc)
	bVersion := Uint642Byte((*index).Version)
	bNUmber := Uint642Byte((*index).Number)

	buffer = append(buffer, bDirc...)
	buffer = append(buffer, bVersion...)
	buffer = append(buffer, bNUmber...)

	for _, entry := range (*index).Entries {

		cUnix := entry.cTime.Unix()
		bcUnix := make([]byte, 4)
		binary.BigEndian.PutUint32(bcUnix, uint32(cUnix))
		buffer = append(buffer, bcUnix...)
		buffer = append(buffer, bcUnix...)

		mUnix := entry.mTime.Unix()
		bmUnix := make([]byte, 4)
		binary.BigEndian.PutUint32(bmUnix, uint32(mUnix))
		buffer = append(buffer, bmUnix...)
		buffer = append(buffer, bmUnix...)

		bDev := make([]byte, 4)
		binary.BigEndian.PutUint32(bDev, entry.Dev)
		buffer = append(buffer, bDev...)

		bInode := make([]byte, 4)
		binary.BigEndian.PutUint32(bInode, entry.Inode)
		buffer = append(buffer, bInode...)

		bMode := make([]byte, 4)
		binary.BigEndian.PutUint32(bMode, entry.Mode)
		buffer = append(buffer, bMode...)

		bUid := make([]byte, 4)
		binary.BigEndian.PutUint32(bUid, entry.Uid)
		buffer = append(buffer, bUid...)

		bGid := make([]byte, 4)
		binary.BigEndian.PutUint32(bGid, entry.Gid)
		buffer = append(buffer, bGid...)

		bSize := make([]byte, 4)
		binary.BigEndian.PutUint32(bSize, entry.Size)
		buffer = append(buffer, bSize...)

		bHash, err := hex.DecodeString(entry.Hash)
		if err != nil {
			continue
		}
		buffer = append(buffer, bHash...)

		bnSize := make([]byte, 2)
		binary.BigEndian.PutUint16(bnSize, uint16(len(entry.Name)))
		buffer = append(buffer, bnSize...)

		bName := []byte(entry.Name)
		buffer = append(buffer, bName...)

		padding := GetPaddingSize(uint64(len(bName)))
		bPadding := make([]byte, padding)
		buffer = append(buffer, bPadding...)
	}

	w, err := os.Create(file_path)
	if err != nil {
		return err
	}
	defer w.Close()

	count, err := w.Write(buffer)
	if err != nil {
		return err
	}
	fmt.Printf("write %d bytes\n", count)

	return errors.New("None")
}
