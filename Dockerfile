# 使用一个基础的 Golang 镜像
FROM amd64/golang:1.22 as build

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 
    # GOPROXY=https://goproxy.cn

# 设置工作目录
WORKDIR /app

# 复制 Golang 项目的源代码到容器中
COPY . /app

# 在容器中编译 Golang 项目
RUN go build -o hermes main.go

# 创建最终的生产镜像
FROM alpine:latest as prod

# 设置工作目录
WORKDIR /app

# 从之前的阶段复制二进制文件
COPY --from=build /app/conf/* /app/conf/
COPY --from=build /app/hermes /app/

# 设置环境变量等
ENV PORT=19876
ENV GIN_MODE=release
ENV HERMES_DATA_DIR=/app/data/hermes

# 暴露端口
EXPOSE 19876

# 启动应用
CMD ["./hermes"]
