# ========== 阶段1：构建阶段（一次性编译两个main.go，生成两个二进制）==========
FROM golang:1.25-alpine AS builder
# 编译环境配置（静态编译+国内加速，杜绝依赖问题）
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct
WORKDIR /app

# 提前下载依赖（利用Docker缓存，大幅提速构建）
COPY go.mod go.sum ./
RUN go mod download

# 复制整个项目源码（包含worker/、watch/两个目录）
COPY . .

# ✅ 编译1：worker/main.go → 输出二进制文件 /bin/worker
RUN go build -ldflags="-s -w" -o /bin/worker worker/main.go

# ✅ 编译2：watch/main.go → 输出二进制文件 /bin/watch
RUN go build -ldflags="-s -w" -o /bin/watch watch/main.go

# ========== 阶段2：运行阶段（复制两个二进制+赋予权限+时区配置）==========
FROM alpine:3.19 AS runner
# 时区配置（解决Go程序时区偏差，上海时区）
RUN apk add --no-cache tzdata && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
WORKDIR /app

# ✅ 复制编译好的两个二进制文件到运行镜像
COPY --from=builder /bin/worker ./
COPY --from=builder /bin/watch ./

# ✅ 必加：给两个二进制文件赋予【执行权限】（彻底解决not found报错）
RUN chmod +x ./worker ./watch