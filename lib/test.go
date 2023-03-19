package lib



import (
	"fmt"
	"encoding/binary"
	"encoding/hex"

	// "reflect"
)

/** インデックスへの書き込みの際のバイト列にバグが発生 */

type Old struct {
	Name string
	Buffer []byte
}


func TestIndex() {
	indexPath := ".git/index"

	oldBuffer, err := GetIndexFile(indexPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	index, _ := CreateIndex(oldBuffer)




	numOfEntry := Bytes2Uint32(oldBuffer[8:12])
	oldBuffer = oldBuffer[12:]

	var olds []Old


	var count uint32
	count = 0
	for {
		if count >= numOfEntry {
			break
		}
		count++

		nsize := Bytes2uint(oldBuffer[60:62])
		name := string(oldBuffer[62 : 62+nsize])
		fmt.Println(name)
		padding := GetPaddingSize(62+nsize)
		offset := 62+nsize+padding

		fmt.Println("padding->", padding, ": offset:>", offset)

		old := Old{
			Name: name,
			Buffer: oldBuffer[:offset],
		}

		olds = append(olds, old)

		oldBuffer = oldBuffer[offset:]
	}

	fmt.Println("len->", len(olds), "NOE->", numOfEntry)



	buffer := make([]byte, 0)

	// bDirc := []byte((*index).Dirc)
	// bVersion := Uint642Byte((*index).Version)
	// bNUmber := Uint642Byte((*index).Number)

	// buffer = append(buffer, bDirc...)
	// buffer = append(buffer, bVersion...)
	// buffer = append(buffer, bNUmber...)

	for _, entry := range (*index).Entries {

		// fmt.Println(entry)
		buffer = make([]byte, 0)

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

		var sw uint64
		sw = 62

		padding := GetPaddingSize(sw+uint64(len(bName)))
		bPadding := make([]byte, padding)

		buffer = append(buffer, bPadding...)
		fmt.Println("")
		fmt.Println("")
		for _, old := range olds {
			if old.Name == entry.Name {
				fmt.Println("OLD->", old.Buffer)
			}
		}
		fmt.Println("NEW->", buffer)
	}

	// fmt.Println(reflect.DeepEqual(oldBuffer, buffer))


}