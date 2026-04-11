export PATH := $(PATH):$(shell go env GOPATH)/bin

gen:
	rm -rf internal/services/common/gen
	protoc \
		--proto_path=protobuf \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=module=github.com/1kyryll/go-grpc \
		--go-grpc_opt=module=github.com/1kyryll/go-grpc \
		protobuf/*.proto
