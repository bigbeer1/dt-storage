.PHONY: build
build:  archive


# 日志微服务
archive:
	go build -o deploy/app/td-st dt-storage/td-storage