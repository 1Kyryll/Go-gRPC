export PATH := $(PATH):$(shell go env GOPATH)/bin

gen:
	protoc \
		--proto_path=protobuf \
		--go_out=services/common/orders \
		--go-grpc_out=services/common/orders \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		protobuf/*.proto
