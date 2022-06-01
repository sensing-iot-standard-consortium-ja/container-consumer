package Common

import (
	"regexp"
	"strconv"
)

func GetDataIndexLength(dataIndex uint8) int {
	dataIdLength := 0
	switch dataIndex {
	case 0x00:
		dataIdLength = 16
		break
	default:
		dataIdLength = 0
		break
	}
	return dataIdLength
}

func SerializeDataId(s string) []byte {
	rex := regexp.MustCompile("^0x")
	dataIdString := rex.ReplaceAllString(s, "")
	dataId := make([]byte, len(dataIdString)/2)
	for i := 0; i < len(dataId); i++ {
		hexHigher, _ := strconv.ParseInt(string(dataIdString[2*i+0]), 16, 8)
		hexLower, _ := strconv.ParseInt(string(dataIdString[2*i+1]), 16, 8)
		hex := (hexHigher << 4) | hexLower
		dataId[i] = byte(hex)
	}
	return dataId
}
