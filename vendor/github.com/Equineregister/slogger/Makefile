.PHONY: lint
lint:
	golangci-lint --enable gosec,misspell run ./...

run-ex: 
ifeq ($(LOGGER_STYLE),)
	@echo Please set the LOGGER_STYLE variable to the build the specific example package you want to test.
	@echo For example: make build-example LOGGER_STYLE=json 
else 
	go run github.com/Equineregister/slogger/example/$(LOGGER_STYLE)
endif
