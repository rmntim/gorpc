all: generate

generate:
	protoc -I contracts/proto contracts/proto/sso/*.proto --go_out=./contracts/gen/go/ --go_opt=paths=source_relative --go-grpc_out=./contracts/gen/go/ --go-grpc_opt=paths=source_relative
