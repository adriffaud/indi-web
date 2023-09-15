BINARY_NAME = indi-web
AIR_URL := https://github.com/cosmtrek/air/releases/latest/download/air-linux-amd64

.PHONY: all
all: build

.PHONY: build
build:
	go build -o $(BINARY_NAME) ./cmd/indi-web

.PHONY: download
download:
	@if [ ! -f ./bin/air ]; then \
		curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s; \
	fi

.PHONY: run
run: download
	./bin/air
