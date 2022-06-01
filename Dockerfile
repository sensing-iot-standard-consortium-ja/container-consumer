FROM golang:1.18.2-buster as BUILDENV
COPY . /app
WORKDIR /app
RUN go build

FROM alpine
COPY --from=BUILDENV /app/hello /
ENTRYPOINT ["/hello"]