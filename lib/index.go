package lib

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"
	"time"

	// "reflect"
	"encoding/hex"
	// "io/ioutil"
)

type Entry struct {
	CTime time.Time
	MTime time.Time
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

const PADDING_BYTE_TO_UINT64_IS_8 = 8
const PADDING_BYTE_TO_UINT32_IS_4 = 4
const PADDING_BYTE_TO_UINT16_IS_2 = 2

func Bytes2uint64(bytes []byte) uint64 {
	// byte列を, entry objectに変換する際に利用
	padding := make([]byte, PADDING_BYTE_TO_UINT64_IS_8-len(bytes))
	source := append(padding, bytes...)
	artifact := binary.BigEndian.Uint64(source)
	return artifact
}

func Bytes2Uint32(bytes []byte) uint32 {
	padding := make([]byte, PADDING_BYTE_TO_UINT32_IS_4-len(bytes))
	source := append(padding, bytes...)
	artifact := binary.BigEndian.Uint32(source)
	return artifact
}

func Bytes2Uint16(bytes []byte) uint16 {
	padding := make([]byte, PADDING_BYTE_TO_UINT16_IS_2-len(bytes))
	source := append(padding, bytes...)
	artifact := binary.BigEndian.Uint16(source)
	return artifact
}

func Byte2UnixTimeStamp(times uint64) (time.Time, error) {
	int64time := int64(times)
	var offsetHour, offsetMinute int
	if _, err := fmt.Sscanf("+0900", "+%02d%02d", &offsetHour, &offsetMinute); err != nil {
		return time.Time{}, err
	}
	location := time.FixedZone(" ", 3600*offsetHour+60*offsetMinute)
	timestamp := time.Unix(int64time, 0).In(location)
	// time.Now().String()
	return timestamp, nil
}

func Byte2Mode(bytes []byte) (uint32, error) {
	dec := binary.BigEndian.Uint32(bytes)
	oct := fmt.Sprintf("%o", dec)
	num, err := strconv.ParseUint(oct, 10, 32)
	if err != nil {
		return 0, err
	}
	uint32Num := uint32(num)
	return uint32Num, nil
}

func Element2byte32(index_line uint32) []byte {
	// entry objectをバイト列に変換する際に利用
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, index_line)
	return buffer
}

func GetPaddingSize(had uint64) uint64 {
	Rem := had % 8
	return 8 - Rem
}

func (c *Client) GetIndexObject(file_path string) (Index, error) {
	f, err := os.Open(file_path)
	if err != nil {
		return Index{}, err
	}
	defer f.Close()

	buffer, err := ioutil.ReadAll(f)
	if err != nil {
		return Index{}, err
	}

	dirc := string(buffer[0:4])
	if dirc != "DIRC" {
		return Index{}, err
	}

	version := Bytes2Uint32(buffer[4:8])
	if version != 2 {
		return Index{}, err
	}

	number_of_entry := Bytes2Uint32(buffer[8:12])

	buffer = buffer[12:]

	var index Index

	index.Dirc = dirc
	index.Version = version
	index.Number = number_of_entry

	var count uint32
	count = 0
	for {
		fmt.Println()
		if count >= number_of_entry {
			break
		}
		count++

		c_time_64 := Bytes2uint64(buffer[0:4])
		c_time, err := Byte2UnixTimeStamp(c_time_64)
		// CTime, err := util.GetUnixTime(util.Bytes2uint(buffer[0:4]))
		if err != nil {
			fmt.Println(err)
			continue
		}
		// _ = Bytes2str(buffer[4:8])

		m_time_64 := Bytes2uint64(buffer[8:12])
		m_time, err := Byte2UnixTimeStamp(m_time_64)
		// MTime, err := GetUnixTime(Bytes2uint(buffer[8:12]))
		if err != nil {
			fmt.Println(err)
			continue
		}
		// _ = Bytes2str(buffer[12:16])

		dev := Bytes2Uint32(buffer[16:20])
		inode := Bytes2Uint32(buffer[20:24])
		mode, err := Byte2Mode(buffer[24:28])
		if err != nil {
			continue
		}
		uid := Bytes2Uint32(buffer[28:32])
		gid := Bytes2Uint32(buffer[32:36])
		size := Bytes2Uint32(buffer[36:40])

		hash := hex.EncodeToString(buffer[40:60])

		nsize := Bytes2uint64(buffer[60:62])
		name := string(buffer[62 : 62+nsize])

		entry := Entry{
			CTime: c_time,
			MTime: m_time,
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
	return index, nil
}

func InitIndexObject() Index {
	var index Index
	index.Dirc = "DIRC"
	index.Version = 2
	index.Number = uint32(0)

	return index
}

func (index *Index) UpdateIndex(name string, hash string) Index {
	current_dir, _ := os.Getwd()
	file_path := current_dir + "/" + name

	var system_call syscall.Stat_t
	syscall.Stat(file_path, &system_call)

	file_info, err := os.Stat(file_path)
	if err != nil {
		fmt.Println(file_path, err)
	}

	oct := fmt.Sprintf("%o", uint32(system_call.Mode))
	mode_number, err := strconv.ParseUint(oct, 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	mode := uint32(mode_number)

	new_entry := Entry{
		CTime: file_info.ModTime(),
		MTime: file_info.ModTime(),
		Dev:   uint32(system_call.Dev),
		Inode: uint32(system_call.Ino),
		Mode:  mode,
		Uid:   system_call.Uid,
		Gid:   system_call.Gid,
		Size:  uint32(system_call.Size),
		Hash:  hash,
		Name:  name,
	}

	var new_index Index

	for _, entry := range index.Entries {
		if entry.Name == name {
			continue
		}
		if entry.Hash == hash {
			continue
		}
		new_index.Entries = append(new_index.Entries, entry)
	}
	new_index.Entries = append(new_index.Entries, new_entry)

	new_index.Dirc = "DIRC"
	new_index.Version = 2
	new_index.Number = uint32(len(new_index.Entries))

	return new_index
}

func (index *Index) AsByte() GitBuffer {
	buffer := make([]byte, 0)
	dirc := []byte(index.Dirc)
	version := Element2byte32(index.Version)
	number := Element2byte32(index.Number)

	buffer = append(buffer, dirc...)
	buffer = append(buffer, version...)
	buffer = append(buffer, number...)

	for _, entry := range index.Entries {

		c_unix := entry.CTime.Unix()
		buf_c_unix := Element2byte32(uint32(c_unix))
		buffer = append(buffer, buf_c_unix...)
		buffer = append(buffer, buf_c_unix...)

		m_unix := entry.MTime.Unix()
		buf_m_unix := Element2byte32(uint32(m_unix))
		buffer = append(buffer, buf_m_unix...)
		buffer = append(buffer, buf_m_unix...)

		dev := entry.Dev
		buf_dev := Element2byte32(uint32(dev))
		buffer = append(buffer, buf_dev...)

		inode := entry.Inode
		buffer_inode := Element2byte32(uint32(inode))
		buffer = append(buffer, buffer_inode...)

		mode := entry.Mode
		buffer_mode := Element2byte32(uint32(mode))
		buffer = append(buffer, buffer_mode...)

		uid := entry.Uid
		buffer_uid := Element2byte32(uint32(uid))
		buffer = append(buffer, buffer_uid...)

		gid := entry.Gid
		buffer_gid := Element2byte32(uint32(gid))
		buffer = append(buffer, buffer_gid...)

		size := entry.Size
		buffer_size := Element2byte32(uint32(size))
		buffer = append(buffer, buffer_size...)

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

		var sw uint64
		sw = 62

		padding := GetPaddingSize(sw + uint64(len(bName)))
		bPadding := make([]byte, padding)
		buffer = append(buffer, bPadding...)
	}

	return GitBuffer{Buffer: buffer}

}

// func UpdateIndex(index *Index, name string, hash string, client *Client) (*Index, string, error) {
// 	current, _ := os.Getwd()
// 	filePath := current + "/" + name

// 	var sysC syscall.Stat_t
// 	syscall.Stat(filePath, &sysC)

// 	fileInfo, err := os.Stat(filePath)
// 	if err != nil {
// 		fmt.Println(filePath)
// 		return nil, "", err
// 	}
// 	oct := fmt.Sprintf("%o", uint32(sysC.Mode))
// 	num, err := strconv.ParseUint(oct, 10, 32)
// 	if err != nil {
// 		return nil, "", err
// 	}
// 	mode := uint32(num)
// 	Nentry := Entry{
// 		CTime: fileInfo.ModTime(),
// 		MTime: fileInfo.ModTime(),
// 		Dev:   uint32(sysC.Dev),
// 		Inode: uint32(sysC.Ino),
// 		Mode:  mode,
// 		Uid:   sysC.Uid,
// 		Gid:   sysC.Gid,
// 		Size:  uint32(sysC.Size),
// 		Hash:  hash,
// 		Name:  name,
// 	}

// 	fmt.Printf("%+v\n", Nentry)

// 	var Nindex Index

// 	Nindex.Dirc = (*index).Dirc
// 	Nindex.Version = (*index).Version
// 	Nindex.Number = (*index).Number

// 	for _, entry := range (*index).Entries {
// 		if entry.Name == name {
// 			continue
// 		}
// 		if entry.Hash == hash {
// 			continue
// 		}
// 		Nindex.Entries = append(Nindex.Entries, entry)
// 		Nentry.CTime = entry.CTime /** ここでファイルの作成時間を遺伝*/
// 	}

// 	Nindex.Entries = append(Nindex.Entries, Nentry)

// 	return &Nindex, filePath, nil
// }

// func Uint642Byte(ui uint32) []byte {
// 	bt := make([]byte, 4)
// 	binary.BigEndian.PutUint32(bt, ui)
// 	return bt
// }

// func WriteIndex(index *Index, file_path string) error {
// 	buffer := make([]byte, 0)

// 	bDirc := []byte((*index).Dirc)
// 	bVersion := Uint642Byte((*index).Version)
// 	bNUmber := Uint642Byte((*index).Number)

// 	buffer = append(buffer, bDirc...)
// 	buffer = append(buffer, bVersion...)
// 	buffer = append(buffer, bNUmber...)

// 	for _, entry := range (*index).Entries {

// 		fmt.Println(entry)

// 		cUnix := entry.CTime.Unix()
// 		bcUnix := make([]byte, 4)
// 		binary.BigEndian.PutUint32(bcUnix, uint32(cUnix))
// 		buffer = append(buffer, bcUnix...)
// 		buffer = append(buffer, bcUnix...)

// 		mUnix := entry.MTime.Unix()
// 		bmUnix := make([]byte, 4)
// 		binary.BigEndian.PutUint32(bmUnix, uint32(mUnix))
// 		buffer = append(buffer, bmUnix...)
// 		buffer = append(buffer, bmUnix...)

// 		bDev := make([]byte, 4)
// 		binary.BigEndian.PutUint32(bDev, entry.Dev)
// 		buffer = append(buffer, bDev...)

// 		bInode := make([]byte, 4)
// 		binary.BigEndian.PutUint32(bInode, entry.Inode)
// 		buffer = append(buffer, bInode...)

// 		bMode := make([]byte, 4)
// 		binary.BigEndian.PutUint32(bMode, entry.Mode)
// 		buffer = append(buffer, bMode...)

// 		bUid := make([]byte, 4)
// 		binary.BigEndian.PutUint32(bUid, entry.Uid)
// 		buffer = append(buffer, bUid...)

// 		bGid := make([]byte, 4)
// 		binary.BigEndian.PutUint32(bGid, entry.Gid)
// 		buffer = append(buffer, bGid...)

// 		bSize := make([]byte, 4)
// 		binary.BigEndian.PutUint32(bSize, entry.Size)
// 		buffer = append(buffer, bSize...)

// 		bHash, err := hex.DecodeString(entry.Hash)
// 		if err != nil {
// 			continue
// 		}
// 		buffer = append(buffer, bHash...)

// 		bnSize := make([]byte, 2)
// 		binary.BigEndian.PutUint16(bnSize, uint16(len(entry.Name)))
// 		buffer = append(buffer, bnSize...)

// 		bName := []byte(entry.Name)
// 		buffer = append(buffer, bName...)

// 		var sw uint64
// 		sw = 62

// 		padding := GetPaddingSize(sw + uint64(len(bName)))
// 		bPadding := make([]byte, padding)
// 		buffer = append(buffer, bPadding...)
// 	}

// 	w, err := os.Create(file_path)
// 	if err != nil {
// 		return err
// 	}
// 	defer w.Close()

// 	count, err := w.Write(buffer)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Printf("write %d bytes\n", count)

// 	return errors.New("None")
// }
