# build host controller binary
FROM golang:1.18 as builder
WORKDIR /go/src/github.com/shani1998/k8s-utility-controller
COPY . .
RUN  CGO_ENABLED=0 go build -mod=vendor -o bin/server ./cmd/

# copy binary from builder
FROM alpine:latest as runner
RUN apk --no-cache add curl
WORKDIR /bin
COPY --from=builder ["/go/src/github.com/shani1998/k8s-utility-controller/bin", "./"]
ENTRYPOINT ["server"]