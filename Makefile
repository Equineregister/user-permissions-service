OUTPUT := bin/service

.PHONY: test
test: lint
	go test -v -tags test ./... --count=1 --race

.PHONY: lint
lint:
	golangci-lint run --modules-download-mode=vendor --enable=gosec,misspell ./...

.PHONY: swagger
swagger:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag fmt
	swag init --dir internal/adapters/primary/rest --generalInfo server.go --parseInternal --output api --outputTypes yaml,json

.PHONY: build
build: lint
	go build -o $(OUTPUT) ./cmd/server/
	$(warning *** You must run 'make swagger' if you need to update your REST API's swagger documentation ***)

.PHONY: run
run: build-debug
	./run.sh

.PHONY: build-debug
build-debug: lint
	go build --race -gcflags="all=-N -l" -ldflags=-linkmode=internal -o $(OUTPUT) ./cmd/server/

.PHONY: govulncheck
govulncheck:
	govulncheck ./...

.PHONY: proto
proto:
	protoc \
	--go_out=. \
	--go_opt=paths=source_relative \
    --go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
    pkg/proto/*.proto
