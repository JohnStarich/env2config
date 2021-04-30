.PHONY: lint
lint:
	go vet ./...

.PHONY: test
test:
	go test -race -cover ./...

.PHONY: build
build:
	docker build -t johnstarich/env2config:latest .

.PHONY: docs
docs:
	go install github.com/johnstarich/go/gopages@v0.1
	gopages \
		-gh-pages \
		-gh-pages-user "${GIT_USER}" \
		-gh-pages-token "${GIT_TOKEN}"
