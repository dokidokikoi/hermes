FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend
COPY izumi-front/ ./
RUN npm install -g pnpm
RUN pnpm install
RUN pnpm run build

FROM golang:1.24-alpine AS backend-builder

RUN sed -i 's#https\?://dl-cdn.alpinelinux.org/alpine#https://mirrors.tuna.tsinghua.edu.cn/alpine#g' /etc/apk/repositories && \ 
    apk update
RUN apk add --no-cache build-base sqlite-dev
WORKDIR /app/backend
COPY ./ ./
RUN go env -w GOPROXY='https://goproxy.cn,direct' && go mod tidy && CGO_ENABLED=1 go build -ldflags="-linkmode external -extldflags '-static'" -o server .

FROM nginx:1.27-alpine

# 拷贝前端静态资源
COPY --from=frontend-builder /app/frontend/dist /usr/share/nginx/html

# 拷贝 Nginx 配置
COPY nginx/default.conf /etc/nginx/conf.d/default.conf

WORKDIR /app

# 拷贝 Go 后端可执行文件
COPY --from=backend-builder /app/backend/server ./

RUN mkdir -p /app/data

# 设置环境变量
ENV GIN_MODE=release
ENV HERMES_DATA_DIR=/data/izimu

EXPOSE 80

CMD ["/bin/sh", "-c", "./server & nginx -g 'daemon off;'"]
