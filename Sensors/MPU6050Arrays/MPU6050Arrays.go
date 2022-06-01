package MPU6050Arrays

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"protoschema/Common"
	"protoschema/Container"
	"reflect"

	"protoschema/Sensors/MPU6050"
)

const DataIndex uint8 = 0x00
const DataIdString string = "0xDEADBEEF1001AAAA0000000000000006"

var DataId []byte

func init() {
	DataId = Common.SerializeDataId(DataIdString)
}

type MPU6050Arrays struct {
	Accel [3]uint16 `json:"Accel"`
	Temp  uint16    `json:"Temp"`
	Gyro  [3]uint16 `json:"Gyro"`
}

func New() *MPU6050Arrays {
	return &MPU6050Arrays{
		Accel: [3]uint16{0x0000, 0x0000, 0x0000},
		Temp:  0x0000,
		Gyro:  [3]uint16{0x0000, 0x0000, 0x0000},
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

func Unmarshal(m *MPU6050Arrays) []byte {
	buf := &bytes.Buffer{}
	_ = binary.Write(buf, binary.BigEndian, m)
	return buf.Bytes()
}

func Marshal(byteArray []byte) *MPU6050Arrays {
	buf := bytes.NewReader(byteArray)
	m := MPU6050Arrays{}
	if err := binary.Read(buf, binary.BigEndian, &m); err != nil {
		fmt.Println("binary.Read failed:", err)
	}

	return &m
}

func (m MPU6050Arrays) Print() {
	fmt.Printf("Payload Format<MPU6050Arrays>\n")
	fmt.Printf("  Accel:     %#v\n", m.Accel)
	fmt.Printf("  Temp:      0x%02X\n", m.Temp)
	fmt.Printf("  Gyro:      %#v\n", m.Gyro)
}

// MPU6050型をMPU6050Arrays型に変換する関数
func ConvertMPU6050toMPU6050Arrays(mpu6050 *MPU6050.MPU6050, mpu6050Arrays *MPU6050Arrays) {
	mpu6050Arrays.Accel = [3]uint16{mpu6050.AccelX, mpu6050.AccelY, mpu6050.AccelZ}
	mpu6050Arrays.Temp = mpu6050.Temp
	mpu6050Arrays.Gyro = [3]uint16{mpu6050.GyroX, mpu6050.GyroY, mpu6050.GyroZ}
}
