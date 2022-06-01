package MPU6050

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"protoschema/Common"
	"protoschema/Container"
	"reflect"
)

const DataIndex uint8 = 0x00
const DataIdString string = "0xDEADBEEF0001AAAA0000000000000006"

var DataId []byte

func init() {
	DataId = Common.SerializeDataId(DataIdString)
}

// strconv.ParseInt(DataIdString, 16, 8)

type MPU6050 struct {
	AccelX uint16 `json:"AccelX"`
	AccelY uint16 `json:"AccelY"`
	AccelZ uint16 `json:"AccelZ"`
	Temp   uint16 `json:"Temp"`
	GyroX  uint16 `json:"GyroX"`
	GyroY  uint16 `json:"GyroY"`
	GyroZ  uint16 `json:"GyroZ"`
}

func New() *MPU6050 {
	return &MPU6050{
		AccelX: 0x0000,
		AccelY: 0x0000,
		AccelZ: 0x0000,
		Temp:   0x0000,
		GyroX:  0x0000,
		GyroY:  0x0000,
		GyroZ:  0x0000,
	}
}

func NewContainer() *Container.Container {
	c := Container.New()
	p := New()
	b := Unmarshal(p)

	c.Header.Type = 0x5555

	dataIndexLength := Common.GetDataIndexLength(c.Header.DataIndex)
	c.Header.Length = uint16(int(reflect.TypeOf(c.Header.Type).Size()) +
		int(reflect.TypeOf(c.Header.Length).Size()) +
		int(reflect.TypeOf(c.Header.DataIndex).Size()) +
		dataIndexLength +
		len(b))

	buf := &bytes.Buffer{}
	_ = binary.Write(buf, binary.BigEndian, DataId)
	c.Header.DataIndex = DataIndex
	c.Header.DataId = buf.Bytes()
	c.Payload = make([]byte, len(b))
	_ = copy(c.Payload, b)

	return c
}

func Unmarshal(m *MPU6050) []byte {
	buf := &bytes.Buffer{}
	_ = binary.Write(buf, binary.BigEndian, m)
	return buf.Bytes()
}

func Marshal(byteArray []byte) *MPU6050 {
	buf := bytes.NewReader(byteArray)
	m := MPU6050{}
	if err := binary.Read(buf, binary.BigEndian, &m); err != nil {
		fmt.Println("binary.Read failed:", err)
	}

	return &m
}

func (m MPU6050) Print() {
	fmt.Printf("Payload Format<MPU6050>\n")
	fmt.Printf("  AccelX:    0x%02X\n", m.AccelX)
	fmt.Printf("  AccelY:    0x%02X\n", m.AccelY)
	fmt.Printf("  AccelZ:    0x%02X\n", m.AccelZ)
	fmt.Printf("  Temp:      0x%02X\n", m.Temp)
	fmt.Printf("  GyroX:     0x%02X\n", m.GyroX)
	fmt.Printf("  GyroY:     0x%02X\n", m.GyroY)
	fmt.Printf("  GyroZ:     0x%02X\n", m.GyroZ)
}
