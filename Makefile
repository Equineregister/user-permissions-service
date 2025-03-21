OUTPUT := bin/service
OUTPUT_LAMBDA_GETUSERPERMS := bin/lambdas/lambda_get_user_permissions
SWAGGER_DIR := internal/adapter/chi
SERVICE := user-permissions-service

.PHONY: test
test: lint
	go test -v -tags test ./... --count=1 --race

.PHONY: lint
lint:
	golangci-lint run --modules-download-mode=vendor --enable=gosec,misspell ./...

.PHONY: compile-lambdas
compile-lambdas:
	go build -o $(OUTPUT_LAMBDA_GETUSERPERMS)/lambda_get_user_permissions github.com/Equineregister/$(SERVICE)/cmd/lambda_get_user_permissions

.PHONY: compile
compile:
	go build -o $(OUTPUT) github.com/Equineregister/$(SERVICE)/cmd/server

.PHONY: build-lambdas
build-lambdas: compile-lambdas lint 

.PHONY: build
build: compile lint 

.PHONY: run
run: build-debug
	./run.sh

.PHONY: build-debug
build-debug: lint 
	go build -race -gcflags="all=-N -l" -ldflags=-linkmode=internal -o $(OUTPUT) github.com/Equineregister/$(SERVICE)/cmd/server

.PHONY: govulncheck
govulncheck:
	govulncheck ./...

.PHONY: swagger-install
swagger-install:
	go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: swagger
swagger:
	swag fmt
	swag init --dir $(SWAGGER_DIR) --generalInfo server.go --parseInternal --output api --outputTypes yaml,json
