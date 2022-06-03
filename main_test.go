package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"protoschema/Container"
	"protoschema/Schema"
	"testing"
)

func TestCreateData(t *testing.T) {
	container_hexascii := "000000490000112233445566778899aabbccddeeff1d17bcbbbfda7efb68a82a40c019f1c5a92eed6440838c6cab4ec4407668fd38adaade4047752df44eb70c02296b8194cbf0cbf0"
	container_bin, _ := hex.DecodeString(container_hexascii)
	f, _ := os.Create("test.bin")
	f.Write(container_bin)
	f.Close()
}

func TestContainer(t *testing.T) {
	container_hexascii := "000000490000112233445566778899aabbccddeeff1d17bcbbbfda7efb68a82a40c019f1c5a92eed6440838c6cab4ec4407668fd38adaade4047752df44eb70c02296b8194cbf0cbf0"
	container_bin, _ := hex.DecodeString(container_hexascii)
	container := Container.Marshal(container_bin)
	// f, _ := os.Create("test.bin")
	// f.Write(container_bin)
	// f.Close()

	if container.Header.Length != 0x49 {
		t.Errorf("container length == 0x49")
	}
	dataIdBin, _ := hex.DecodeString("00112233445566778899aabbccddeeff")

	if !bytes.Equal(container.Header.DataId, dataIdBin) {
		t.Errorf("DataID is 0x00112233445566778899aabbccddeeff")
	}

}
func TestPayloadParse(t *testing.T) {
	container_hexascii := "000000490000112233445566778899aabbccddeeff1d17bcbbbfda7efb68a82a40c019f1c5a92eed6440838c6cab4ec4407668fd38adaade4047752df44eb70c02296b8194cbf0cbf0"
	container_bin, _ := hex.DecodeString(container_hexascii)
	container := Container.Marshal(container_bin)

	dataIdBin, _ := hex.DecodeString("00112233445566778899aabbccddeeff")

	if !bytes.Equal(container.Header.DataId, dataIdBin) {
		t.Errorf("DataID is 0x00112233445566778899aabbccddeeff")
	}
	f, _ := os.Open("tests_examples/0_00112233445566778899aabbccddeeff.json")
	// Schemaの読み込み
	schema, _ := ioutil.ReadAll(f)

	s := Schema.Schema{}
	json.Unmarshal([]byte(schema), &s)
	stp, _ := s.Marshal(container.Payload)
	for _, ss := range stp {
		fmt.Println(ss.Name)
		fmt.Println(ss.Value)
		fmt.Println(ss.Payload)
	}

	// goroutin
	// for field := range s.Fields {
	// 	field.Unmarshal()
	// }

}
