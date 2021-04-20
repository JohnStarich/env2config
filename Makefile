.PHONY: lint
lint:
	go vet ./...

.PHONY: test
test:
	go test -race -cover ./...

.PHONY: build
build:
	docker build -t johnstarich/env2config:latest .
