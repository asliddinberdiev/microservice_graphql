# envs
account_db = postgres://postgres:password@localhost:5432/postgres?sslmode=disable
order_db = postgres://postgres:password@localhost:5432/postgres?sslmode=disable

tidy: 
	@go mod tidy

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

order-proto:
	@protoc \
		--proto_path=order "order/order.proto" \
		--go_out=order/proto --go_opt=paths=source_relative \
		--go-grpc_out=order/proto \
		--go-grpc_opt=paths=source_relative

# generate graphql schema
graphql-schema:
	@cd graphql && go run github.com/99designs/gqlgen generate

# psql migrations
migrate-up:
	@migrate --path ./account/migrations --database "$(account_db)" --verbose up
	@migrate --path ./order/migrations --database "$(order_db)" --verbose up

migrate-down:
	@migrate --path ./account/migrations --database "$(account_db)" --verbose down
	@migrate --path ./order/migrations --database "$(order_db)" --verbose down

# run grpc servers
run-account:
	@go run ./account/cmd/account/main.go

run-catalog:
	@go run ./catalog/cmd/catalog/main.go

run-order:
	@go run ./order/cmd/order/main.go

run-graphql:
	@go run ./graphql/cmd/graphql/main.go

compose-up:
	@docker compose up -d --build

compose-down:
	@docker compose down

compose-logs:
	@docker compose logs -f

compose-restart:
	@docker compose restart