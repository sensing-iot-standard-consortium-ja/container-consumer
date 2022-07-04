package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"protoschema/Container"
	"protoschema/Schema"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

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
	brokerEndpoint := getEnv("KAFKA_BROKER", "localhost:9092")
	fmt.Println(brokerEndpoint)
	consumer, _ := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokerEndpoint,
		"group.id":          "hoge",
		"auto.offset.reset": "smallest",
	})
	producer, _ := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokerEndpoint,
	})

	subscribe_topic_data := getEnv("KAFKA_SUBSCRIBE_TOPIC", "mb_ctopic")
	subscribe_topics := strings.Split(subscribe_topic_data, ",")
	_ = consumer.SubscribeTopics(subscribe_topics, nil)

	produce_topic_prefix := getEnv("KAFKA_PRODUCER_TOPIC_PREFIX", "jsondev_")
	defer consumer.Close()
	run := true

	schemaCache := sync.Map{}
	for run == true {
		ev := consumer.Poll(0)
		switch e := ev.(type) {
		case *kafka.Message:
			jsonBytes := processContainer(e.Value, &schemaCache)
			topic := fmt.Sprintf("%s_%s", produce_topic_prefix, *e.TopicPartition.Topic)
			fmt.Printf("%s\t%x\n", topic, jsonBytes)
			producer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          jsonBytes,
			}, nil)

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
	hostname := getEnv("IOT_SCHEMA_REGISTORY", "http://localhost:30002")
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

func processContainer(buf []byte, schemaCache *sync.Map) []byte {
	// コンテナ・バイト列bufを，コンテナ型に変換
	container := Container.Marshal(buf)
	// [debug] コンテナの中身を確認
	// container.Print()
	dataId := container.Header.DataId
	dataIndex := container.Header.DataIndex
	schemaKey := SchemaKey{dataIndex, dataId}

	schema, ok := schemaCache.Load(schemaKey.String())
	if !ok {
		schema := retriveSchema(schemaKey)
		schemaCache.Store(schemaKey.String(), schema)
	}

	schema_, _ := schema.(Schema.Schema)
	structData, _ := schema_.Marshal(container.Payload)

	item := make(map[string]interface{})
	for _, ss := range structData {
		item[ss.Name] = ss.Value
	}

	a, _ := json.Marshal(item)
	// fmt.Println(string(a))
	return a
}
