export PATH := $(PATH):$(shell go env GOPATH)/bin

gen:
	-rd /s /q internal\services\common\gen 2>nul
	protoc \
		--proto_path=protobuf \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=module=github.com/1kyryll/go-grpc \
		--go-grpc_opt=module=github.com/1kyryll/go-grpc \
		protobuf/*.proto
