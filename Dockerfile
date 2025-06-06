# 构建阶段
FROM golang:1.24.3-alpine AS builder

# 设置容器中的当前工作目录
WORKDIR /build

# 安装依赖（Git 是 Go 模块所需）
RUN apk add --no-cache git

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

# 复制源代码到容器中
COPY . .

# 下载所有依赖。如果 go.mod 和 go.sum 文件没有变化，依赖会被缓存
RUN go mod tidy

# 构建服务器二进制文件
RUN go build -o deploy/app/td-st dt-storage/td-storage


# 最终阶段
FROM alpine:3.18

# 设置容器中的当前工作目录
WORKDIR /app

# 安装所需的依赖（Alpine 镜像基础非常小）
RUN apk add --no-cache ca-certificates

# 从构建阶段复制二进制文件到最终镜像
COPY --from=builder /build/deploy/app/td-st /app/td-st
COPY --from=builder /build/deploy/app/td-storage.yaml /app/td-storage.yaml

# 启动应用程序（服务器）
CMD ["/app/td-st"]
