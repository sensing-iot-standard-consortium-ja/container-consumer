package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"protoschema/Container"
	"protoschema/Sensors/MPU6050"
	"protoschema/Sensors/MPU6050Arrays"
)

//var endian = binary.BigEndian
var Endian = binary.LittleEndian

func containerStubdata() []byte { // 中身(ペイロード)を作る
	c := MPU6050.NewContainer()
	p := MPU6050.New()
	p.AccelX = 0x1122
	p.AccelY = 0x3344
	p.AccelZ = 0x5566
	p.Temp = 0x7788
	p.GyroX = 0x99AA
	p.GyroY = 0xBBCC
	p.GyroZ = 0xDDEE
	// コンテナから見るとPayloadはバイト列
	copy(c.Payload, MPU6050.Unmarshal(p))
	buf := Container.Unmarshal(c)

	// コンテナから見るとPayloadはバイト列
	return buf
}

func sandbox() {
	var err error
	// テスト用のコンテナを作る
	buf := containerStubdata()
	fmt.Printf("Input: %X\n", buf)
	// ここからテストできる

	// コンテナ・バイト列bufを，コンテナ型に変換
	fmt.Println("MPU6050 Container:")
	mpu6050Container := Container.Marshal(buf)
	// コンテナの中身を確認
	mpu6050Container.Print()
	// コンテナのペイロードをパース
	mpu6050Payload := MPU6050.Marshal(mpu6050Container.Payload)
	// ペイロードの中身を確認
	mpu6050Payload.Print()

	// コンテナの内容をJSONで表示
	buf, err = json.MarshalIndent(mpu6050Container, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Container: %v\n", string(buf))
	// ペイロードの内容をJSONで表示
	buf, err = json.MarshalIndent(mpu6050Payload, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Payload: %v\n", string(buf))

	// 参考例1：
	// センサから来たコンテナのMPU6050型のペイロードを，
	// データの種類ごとに配列にまとめたMPU6050Arrays型に変換する

	// 出力先のコンテナを作る
	mpu6050ArraysContainer := MPU6050Arrays.NewContainer()
	mpu6050ArraysPayload := MPU6050Arrays.Marshal(mpu6050ArraysContainer.Payload)

	// 変換の演算を実行
	fmt.Println("ConvertMPU6050toMPU6050Arrays()")
	MPU6050Arrays.ConvertMPU6050toMPU6050Arrays(mpu6050Payload, mpu6050ArraysPayload)

	// 出力先のコンテナのペイロードに結果を格納する
	_ = copy(mpu6050ArraysContainer.Payload, MPU6050Arrays.Unmarshal(mpu6050ArraysPayload))

	// 出力先のコンテナをバイト列に変換して出力する
	fmt.Println("MPU6050Arrays Container:")
	mpu6050ArraysContainer.Print()
	mpu6050ArraysPayload.Print()

	// コンテナの内容をJSONで表示
	buf, err = json.MarshalIndent(mpu6050ArraysContainer, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Container: %v\n", string(buf))
	// ペイロードの内容をJSONで表示
	buf, err = json.MarshalIndent(mpu6050ArraysPayload, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Payload: %v\n", string(buf))

	buf = Container.Unmarshal(mpu6050ArraysContainer)
	fmt.Printf("Output: %X\n", buf)
}
