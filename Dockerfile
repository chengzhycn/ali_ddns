# 基于适当的基础镜像，如Alpine Linux
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 复制程序文件到容器中
COPY ali_ddns /app/ali_ddns

# 安装依赖项，如有需要
# RUN apk add --no-cache <dependency>

# 设置程序环境变量，如有需要
# ENV VARIABLE_NAME=value

# 定义容器启动时要执行的命令
CMD ["/app/ali_ddns"]