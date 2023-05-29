package util

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"time"
)

const PADDING_BYTE_TO_UINT64_IS_8 = 8
const PADDING_BYTE_TO_UINT32_IS_4 = 4
const PADDING_BYTE_TO_UINT16_IS_2 = 2

func Bytes2uint64(bytes []byte) uint64 {
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
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, index_line)
	return buffer
}

func GetPaddingSize(had uint64) uint64 {
	Rem := had % 8
	return 8 - Rem
}
