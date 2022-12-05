.PHONY: vendor
vendor:
	@export GO111MODULE=on; go mod tidy; go mod vendor;unset GO111MODULE

.PHONY: protobuf-gen
protobuf-gen: #	protoc -I$@/${PROTO_DIR} --go_opt=module=${PACKAGE} --go_out=. --go-grpc_opt=module=${PACKAGE} --go-grpc_out=. $@/${PROTO_DIR}/*.proto
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative  proto/*.proto


about: ## Display info related to the build
	@echo "OS: ${OS}"
	@echo "Shell: ${SHELL} ${SHELL_VERSION}"
	@echo "Protoc version: $(shell protoc --version)"
	@echo "Go version: $(shell go version)"
	@echo "Go package: ${PACKAGE}"
	@echo "Openssl version: $(shell openssl version)"
