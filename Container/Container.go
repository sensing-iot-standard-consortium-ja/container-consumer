package Container

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"protoschema/Common"
	"reflect"
)

type OptionalField struct {
	AttributeId    uint8  `json:"Id"`
	Length         uint8  `json:"Length"`
	AttributeValue []byte `json:"Value"`
}

type Header struct {
	Type          uint16          `json:"Type"`
	Length        uint16          `json:"Length"`
	DataIndex     uint8           `json:"DataIndex"`
	DataId        []byte          `json:"-"`
	OptionalField []OptionalField `json:"OptionalField"`
}

type Container struct {
	Header  Header `json:"Header"`
	Payload []byte `json:"Payload"`
}

func New() *Container {
	header := &Header{
		Type:          0x0000,
		Length:        0x0000,
		DataIndex:     0x00,
		DataId:        make([]byte, 0),
		OptionalField: make([]OptionalField, 0),
	}
	return &Container{
		Header:  *header,
		Payload: make([]byte, 0),
	}
}
func (c *Container) MarshalJSON() ([]byte, error) {
	type Alias Container

	return json.Marshal(&struct {
		*Alias
		DataId string `json:"DataId"`
	}{
		Alias:  (*Alias)(c),
		DataId: hex.EncodeToString(c.Header.DataId),
	})
}
func (c *Container) UnmarshalJSON(b []byte) error {
	type Alias Container
	aux := &struct {
		*Alias
		DataId string `json:"DataId"`
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}
	aux.Header.DataId, _ = hex.DecodeString(aux.DataId)
	return nil
}

func Marshal(byteArray []byte) *Container {
	c := New()
	buf := bytes.NewReader(byteArray)

	// ContainerType
	if err := binary.Read(buf, binary.BigEndian, &c.Header.Type); err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	if err := binary.Read(buf, binary.BigEndian, &c.Header.Length); err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	if err := binary.Read(buf, binary.BigEndian, &c.Header.DataIndex); err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	dataIndexLength := Common.GetDataIndexLength(c.Header.DataIndex)
	c.Header.DataId = make([]byte, dataIndexLength)
	if err := binary.Read(buf, binary.BigEndian, &c.Header.DataId); err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	// TODO: OptField設定はまだない

	payloadSize := int(c.Header.Length) -
		(int(reflect.TypeOf(c.Header.Type).Size()) +
			int(reflect.TypeOf(c.Header.Length).Size()) +
			int(reflect.TypeOf(c.Header.DataIndex).Size()) +
			dataIndexLength)
	c.Payload = make([]byte, payloadSize)
	if err := binary.Read(buf, binary.BigEndian, c.Payload); err != nil {
		fmt.Println("binary.Read failed:", err)
	}

	return c
}

func Unmarshal(c *Container) []byte {
	buf := &bytes.Buffer{}
	_ = binary.Write(buf, binary.BigEndian, c.Header.Type)
	_ = binary.Write(buf, binary.BigEndian, c.Header.Length)
	_ = binary.Write(buf, binary.BigEndian, c.Header.DataIndex)
	_ = binary.Write(buf, binary.BigEndian, c.Header.DataId)
	_ = binary.Write(buf, binary.BigEndian, c.Payload)

	return buf.Bytes()
}

func (c *Container) Print() {
	fmt.Printf("Type:      0x%02X\n", c.Header.Type)
	fmt.Printf("Length:    0x%02X\n", c.Header.Length)
	fmt.Printf("DataIndex: 0x%01X\n", c.Header.DataIndex)
	fmt.Printf("DataId:    0x%X\n", c.Header.DataId)
	// TODO: OptField設定はまだない
	for idx, optfield := range c.Header.OptionalField {
		fmt.Printf("Optional Field[%2d]: %v\n", idx, optfield)
	}
	// c.Payload.(reflect.TypeOf(c.Payload))
	// c.Payload.Print()
	buf := &bytes.Buffer{}
	_ = binary.Write(buf, binary.BigEndian, c.Payload)
	fmt.Printf("Payload:   0x%X\n", buf.Bytes())
}
