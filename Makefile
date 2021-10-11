BINARY_NAME=sayit
TEST_DIR=~/.sayit/test

ifeq (run,$(firstword $(MAKECMDGOALS)))
  ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(ARGS):;@:)
endif

build:
	go mod tidy
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME} main.go

.PHONY: run
run:
	./${BINARY_NAME} $(ARGS)

clean:
	go clean
	rm ${TEST_DIR} -rf

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

vet:
	go vet
