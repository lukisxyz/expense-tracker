APP_NAME = main
BUILD_DIR = $(PWD)/tmp

clean:
	rm -rf ./build

critic:
	gocritic check -enableAll ./...

security:
	tmp/gosec ./...

lint:
	tmp/golangci-lint run ./...

test: clean critic security lint
	go test -v -timeout 30s -coverprofile=cover.out -cover ./...
	go tool cover -func=cover.out
build: clean critic security lint
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go

run: build
	tmp/air

setup.air:
	curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b tmp/

setup.gocritic:
	go install -v github.com/go-critic/go-critic/cmd/gocritic@latest

setup.gosec:
	curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b tmp/ v2.16.0

setup.lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b tmp/ v1.54.1

setup: setup.air setup.gocritic setup.lint setup.gosec
