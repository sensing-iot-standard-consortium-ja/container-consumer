module protoschema

go 1.18

replace (
	protoschema/Common => ./Common
	protoschema/Container => ./Container
	protoschema/Schema => ./Schema
	protoschema/Sensors/MPU6050 => ./Sensors/MPU6050
	protoschema/Sensors/MPU6050Arrays => ./Sensors/MPU6050Arrays
)

require github.com/confluentinc/confluent-kafka-go v1.8.2
