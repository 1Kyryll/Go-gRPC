export PATH := $(PATH):$(shell go env GOPATH)/bin

proto:
	-rd /s /q internal\services\common\gen 2>nul
	protoc \
		--proto_path=protobuf \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=module=github.com/1kyryll/go-grpc \
		--go-grpc_opt=module=github.com/1kyryll/go-grpc \
		protobuf/*.proto

sqlc:
	sqlc generate

gqlgen:
	go run github.com/99designs/gqlgen generate

gen: proto sqlc gqlgen
