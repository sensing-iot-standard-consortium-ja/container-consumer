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

	"github.com/confluentinc/confluent-kafka-go/kafka"
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
	fmt.Printf("%X\n", container.Payload)
	// goroutin
	// for field := range s.Fields {
	// 	field.Unmarshal()
	// }

}

func TestContainerToKafkaProduce(t *testing.T) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					// parse dataId
					container := Container.Marshal(ev.Value)
					dataId, _ := hex.DecodeString("00112233445566778899aabbccddeeff")
					var dataIdx uint8 = 0
					if bytes.Equal(container.Header.DataId, dataId) && container.Header.DataIndex == dataIdx {
						f, _ := os.Open("tests_examples/0_00112233445566778899aabbccddeeff.json")
						// Schemaの読み込み
						schema, _ := ioutil.ReadAll(f)
						s := Schema.Schema{}
						json.Unmarshal([]byte(schema), &s)
						stPayload, _ := s.Marshal(container.Payload)
						for _, ss := range stPayload {
							fmt.Println(ss.Name)
							fmt.Println(ss.Value)
							fmt.Println(ss.Payload)
						}
					}
				}
			}
		}
	}()

	container_hexascii := "000000490000112233445566778899aabbccddeeff1d17bcbbbfda7efb68a82a40c019f1c5a92eed6440838c6cab4ec4407668fd38adaade4047752df44eb70c02296b8194cbf0cbf0"
	container_bin, _ := hex.DecodeString(container_hexascii)

	// Produce messages to topic (asynchronously)
	topic := "myTopic"
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          container_bin,
	}, nil)

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)

}
