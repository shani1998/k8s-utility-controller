

vendor:
	@go mod tidy; go mod vendor; unset GO111MODULE

test:
	@go vet ./...
	@go test ./...