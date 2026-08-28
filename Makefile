GIT_COMMIT:=$(shell git describe --dirty --always)

# 项目工具下载
install:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest && \
	go install github.com/evilmartians/lefthook@latest && \
	lefthook install

# 依赖下载
tidy:
	go mod tidy

# 构建
.PHONY: build
build:
	go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn && \
    	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w -X main.REVISION=$(GIT_COMMIT)" -a -o build/mw github.com/1303-yzym/MoonshotWell/cmd

# make version-set TAG=0.2.0
version-set:
	@current=$$(grep 'VERSION' cmd/version.go | awk '{ print $$4 }' | tr -d '"'); \
	next="$(TAG)"; \
	sed -i  "s/$$current/$$next/g" cmd/version.go

# make release VERSION=0.2.0
release:
	git tag -s -m $(VERSION) $(VERSION)
