# 使用 alpine:latest 作为基础镜像
FROM alpine:latest

# 替换镜像源为阿里云，加速 apk 安装
RUN sed -i 's|dl-cdn.alpinelinux.org|mirrors.aliyun.com|g' /etc/apk/repositories

# 安装 tzdata 并配置上海时区（保留 tzdata）
RUN apk add --no-cache tzdata && \
    cp -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 将编译好的二进制文件复制到容器中
COPY myapp /app/myapp

# 赋予执行权限
RUN chmod +x /app/myapp

# 设置工作目录
WORKDIR /app

# 设置容器启动命令
CMD ["/app/myapp"]
