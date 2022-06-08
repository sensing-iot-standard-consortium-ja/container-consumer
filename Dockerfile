FROM golang:1.18.2-buster as BUILDENV
COPY . /app
WORKDIR /app
RUN go build

FROM golang:1.18.2-buster
COPY --from=BUILDENV /app/protoschema /protoschema
ENV KAFKA_BROKER=localhost:9092
ENV KAFKA_SUBSCRIBE_TOPIC=mb_ctopic
ENV KAFKA_PRODUCER_TOPIC=mobile_topic
ENV IOT_SCHEMA_REGISTORY=http://localhost:30002
ENTRYPOINT ["/protoschema"]