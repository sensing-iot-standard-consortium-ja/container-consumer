package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"protoschema/Container"
	"protoschema/Schema"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	// Ctrl+Cなどのシグナルで終了するようにする
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGINT,
		os.Interrupt,
		syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		os.Exit(0)
	}()

	// Consumerを作る
	consumer, _ := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "foo",
		"auto.offset.reset": "smallest",
	})

	_ = consumer.SubscribeTopics([]string{"mb_ctopic"}, nil)
	defer consumer.Close()
	run := true

	schemaCache := sync.Map{}
	for run == true {
		ev := consumer.Poll(0)
		switch e := ev.(type) {
		case *kafka.Message:
			processContainer(e.Value, schemaCache)
		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			run = false
		default:
			//fmt.Printf("Ignored %v\n", e)
		}
	}
}

func retriveSchema(schemaKey SchemaKey) Schema.Schema {
	dataIdstr := hex.EncodeToString(schemaKey.dataId)
	hostname := "http://localhost:30002"
	url := fmt.Sprintf("%s/registry/repo/%d/%s", hostname, schemaKey.dataIndex, dataIdstr)
	fmt.Println(url)
	resp, _ := http.Get(url)
	schema_define, _ := ioutil.ReadAll(resp.Body)
	// Schema読み込み
	schema := Schema.Schema{}
	json.Unmarshal(schema_define, &schema)
	return schema
}

type SchemaKey struct {
	dataIndex uint8
	dataId    []byte
}

func (schema *SchemaKey) String() string {
	return fmt.Sprintf("%d_%s", schema.dataIndex, schema.dataId)
}

func processContainer(buf []byte, schemaCache sync.Map) {
	// コンテナ・バイト列bufを，コンテナ型に変換
	fmt.Println("Container:")
	container := Container.Marshal(buf)
	// [debug] コンテナの中身を確認
	// container.Print()
	dataId := container.Header.DataId
	dataIndex := container.Header.DataIndex
	schemaKey := SchemaKey{dataIndex, dataId}

	schema, _ := schemaCache.LoadOrStore(schemaKey.String(), retriveSchema(schemaKey))
	schema_, _ := schema.(Schema.Schema)
	structData, _ := schema_.Marshal(container.Payload)

	for _, ss := range structData {
		fmt.Println(ss.Name)
		fmt.Println(ss.Value)
		fmt.Println(ss.Payload)
	}

}
