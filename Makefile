EXCLUDE_LINT = "_test.go"
GO_TEST_ARGS := -short -p=1 -tags vartime

test_fmt:
	@echo Checking correct formatting of files
	@{ \
		files=$$( go fmt ./... ); \
		if [ -n "$$files" ]; then \
		echo "Files not properly formatted: $$files"; \
		exit 1; \
		fi; \
		if ! go vet ./...; then \
		exit 1; \
		fi \
	}

test_lint:
	@echo Checking linting of files
	@{ \
		GO111MODULE=off go get -u golang.org/x/lint/golint; \
		el=$(EXCLUDE_LINT); \
		lintfiles=$$( golint ./... | egrep -v "$$el" ); \
		if [ -n "$$lintfiles" ]; then \
		echo "Lint errors:"; \
		echo "$$lintfiles"; \
		exit 1; \
		fi \
	}

test_local:
	go test -v $(GO_TEST_ARGS) ./...

test_codecov:
	./coveralls.sh $(GO_TEST_ARGS)

test: test_fmt test_lint test_codecov

local: test_fmt test_lint test_local
