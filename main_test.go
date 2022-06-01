package main

import (
	"encoding/hex"
	"protoschema/Container"
	"testing"
)

func TestContainer(t *testing.T) {
	container_hexascii := "000000490000112233445566778899aabbccddeeff1d17bcbbbfda7efb68a82a40c019f1c5a92eed6440838c6cab4ec4407668fd38adaade4047752df44eb70c02296b8194cbf0"
	container_bin, _ := hex.DecodeString(container_hexascii)
	container := Container.Marshal(container_bin)
	if container.Header.Length != 0x49 {
		t.Errorf("container length == 0x49")
	}

}
