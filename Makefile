# Project Name
SHA1				:= $(shell git rev-parse --verify --short HEAD)
INTERNAL_BUILD_ID	:= $(shell [ -z "${BUILD_ID}" ] && echo "local" || echo ${BUILD_ID})
VERSION				:= $(shell echo "${INTERNAL_BUILD_ID}_${SHA1}")
BINARY				:= $(shell basename -s .git `git config --get remote.origin.url`)


.PHONY: build
build:

	docker run  -t -v $(PWD):/go/src/github.com/jaimemartinez88/$(BINARY) -w /go/src/github.com/jaimemartinez88/$(BINARY) golang:1.10.2 go build -x -ldflags "-X main.version=$(VERSION)" -o $(BINARY) ./cmd/$(BINARY)/

.PHONY: publish
publish:
	echo "publishing"
ifndef ENV
	$(error ENV is not passed [corp|dev|prod])
endif

.PHONY: config
config:
	$(shell go build ./cmd/$(BINARY) && ./$(BINARY) -default config.json)
