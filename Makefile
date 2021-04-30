LINT_VERSION=1.39.0

.PHONY: lint-deps
lint-deps:
	@if ! which golangci-lint >/dev/null || [[ "$$(golangci-lint version 2>&1)" != *${LINT_VERSION}* ]]; then \
		curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v${LINT_VERSION}; \
	fi

.PHONY: lint
lint: lint-deps
	golangci-lint run

.PHONY: lint-fix
lint-fix: lint-deps
	golangci-lint run --fix

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
