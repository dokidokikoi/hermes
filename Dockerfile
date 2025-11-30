FROM docker.m.daocloud.io/node:20-alpine AS frontend-builder

WORKDIR /app/frontend
RUN sed -i 's#https\?://dl-cdn.alpinelinux.org/alpine#https://mirrors.tuna.tsinghua.edu.cn/alpine#g' /etc/apk/repositories && \ 
    apk update && apk add --no-cache git
RUN git clone https://github.com/dokidokikoi/izumi-front.git ./
RUN npm install -g pnpm
RUN pnpm install
RUN pnpm run build

FROM docker.m.daocloud.io/golang:1.24-alpine AS backend-builder

RUN sed -i 's#https\?://dl-cdn.alpinelinux.org/alpine#https://mirrors.tuna.tsinghua.edu.cn/alpine#g' /etc/apk/repositories && \ 
    apk update && apk add --no-cache build-base sqlite-dev
WORKDIR /app/backend
COPY ./ ./
RUN go env -w GOPROXY='https://goproxy.cn,direct' && go mod tidy && CGO_ENABLED=1 go build -ldflags="-linkmode external -extldflags '-static'" -o server .

FROM docker.m.daocloud.io/nginx:1.27-alpine

RUN sed -i 's#https\?://dl-cdn.alpinelinux.org/alpine#https://mirrors.tuna.tsinghua.edu.cn/alpine#g' /etc/apk/repositories && \ 
    apk update && apk add --no-cache supervisor libc6-compat

# 拷贝前端静态资源
COPY --from=frontend-builder /app/frontend/dist /usr/share/nginx/html

# 拷贝 Nginx 配置
COPY nginx/default.conf /etc/nginx/conf.d/default.conf

RUN mkdir -p /data/db /app/log && chmod -R 777 /var/cache/nginx /var/run /var/log/nginx /run

WORKDIR /app

# 拷贝 Go 后端可执行文件
COPY --from=backend-builder /app/backend/server ./

# 设置环境变量
ENV GIN_MODE=release
ENV IZUMI_DATA_DIR=/data/izumi

COPY supervisor/supervisor.conf /etc/supervisor.d/izumi.conf

EXPOSE 8888

CMD ["supervisord", "-c", "/etc/supervisor.d/izumi.conf"]
