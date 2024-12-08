tidy: 
	@go mod tidy
	@go mod vendor

# protobuf generate
account-proto:
	@protoc \
		--proto_path=account "account/account.proto" \
		--go_out=account/proto --go_opt=paths=source_relative \
		--go-grpc_out=account/proto \
		--go-grpc_opt=paths=source_relative
catalog-proto:
	@protoc \
		--proto_path=catalog "catalog/catalog.proto" \
		--go_out=catalog/proto --go_opt=paths=source_relative \
		--go-grpc_out=catalog/proto \
		--go-grpc_opt=paths=source_relative


# run grpc servers
run-account:
	@go run ./account/cmd/account/main.go
run-catalog:
	@go run ./catalog/cmd/catalog/main.go

run-graphql:
	@go run ./graphql/main.go