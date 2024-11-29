tidy: 
	@go mod tidy
	@go mod vendor

graphql:
	@go run ./graphql/main.go