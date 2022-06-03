package Schema

import (
	"bytes"
	"encoding/binary"
)

type Field struct {
	Name   string                 `json:"name"`
	Type   string                 `json:"type"`
	Pos    uint                   `json:"pos"`
	Length uint                   `json:"length"`
	Tags   map[string]interface{} `json:"tags"`
}

type StructPayload struct {
	Payload []byte      `json:"-"`
	Name    string      `json:"name"`
	Value   interface{} `json:"value"`
}

func (schema *Schema) Marshal(s []byte) ([]StructPayload, error) {
	items := []StructPayload{}
	for _, field := range schema.Fields {
		structPayload := StructPayload{
			Name:    field.Name,
			Payload: s[field.Pos : field.Pos+field.Length],
		}
		// structPayload.name = field.Name
		// slice := s[field.Pos : field.Pos+field.Length]
		// structPayload.payload = slice

		var byteOrder binary.ByteOrder
		if _, ok := field.Tags["isLittleEndian"]; ok {
			byteOrder = binary.LittleEndian
		} else {
			byteOrder = binary.BigEndian
		}

		buf := bytes.NewReader(structPayload.Payload)
		value := &structPayload.Value
		switch field.Type {
		case "u8":
			var val uint8
			binary.Read(buf, byteOrder, &val)
			*value = val
		case "u16":
			var val uint16
			binary.Read(buf, byteOrder, &val)
			*value = val
		case "u32":
			var val uint32
			binary.Read(buf, byteOrder, &val)
			*value = val
		case "u64":
			var val uint64
			binary.Read(buf, byteOrder, &val)
			*value = val
		case "int8":
			var val int8
			binary.Read(buf, byteOrder, &val)
			*value = val
		case "i16":
			var val int16
			binary.Read(buf, byteOrder, &val)
			*value = val
		case "i32":
			var val int32
			binary.Read(buf, byteOrder, &val)
			*value = val
		case "f32":
			var val float32
			binary.Read(buf, byteOrder, &val)
			*value = val
		case "f64":
			var val float64
			binary.Read(buf, byteOrder, &val)
			*value = val
		case "bytes":
			*value = structPayload.Payload
		}
		items = append(items, structPayload)
	}

	return items, nil
}

type Schema struct {
	Type   string  `json:"type"`
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
}
