IMAGE_REGISTRY  ?= skp123
IMAGE_NAME      := $(IMAGE_REGISTRY)/k8s-utility-controller
IMAGE_VERSION   := v0.0.1


.PHONY: vendor
vendor:
	@export GO111MODULE=on; go mod tidy; go mod vendor;unset GO111MODULE

.PHONY: test
test:
	@go vet ./...
	@go test -v -cover ./...

.PHONY: docker-build
docker-build:
	@docker build -t $(IMAGE_NAME):$(IMAGE_VERSION) .

.PHONY: docker-push
docker-push:
	@docker push $(IMAGE_NAME):$(IMAGE_VERSION)

.PHONY: run-local
run-local:
	@go run ./cmd/*.go
