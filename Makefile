#  go mod vendor;
.PHONY: vendor
vendor:
	@export GO111MODULE=on; go mod tidy; go mod vendor;unset GO111MODULE

.PHONY: test
test:
	@go vet ./...
	@go test ./...
