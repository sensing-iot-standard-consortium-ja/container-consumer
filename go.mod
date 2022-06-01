module protoschema

go 1.18

replace (
    protoschema/Container => ./Container
    protoschema/Common => ./Common
    protoschema/Sensors/MPU6050 => ./Sensors/MPU6050
    protoschema/Sensors/MPU6050Arrays => ./Sensors/MPU6050Arrays
)
