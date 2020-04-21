all: server client

server:
	go build ./cmd/server

client:
	go build ./cmd/client

protoc:
	@echo "Generating Go protobuf files"
	GO111MODULE=on go mod download github.com/gogo/protobuf && GO111MODULE=on go mod vendor
	protoc --gogofast_out=plugins=grpc:./pkg/proto/ -I ${GOPATH}/pkg/mod/github.com/gogo/protobuf@v1.3.1/protobuf -I vendor -I vendor/github.com/gogo/protobuf -I ./pkg/proto connect.proto

clean:
	rm -f ./client ./server

.PHONY: protoc clean