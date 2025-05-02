# 生成wire依赖注入
go tool github.com/google/wire/cmd/wire ./...
# 配置
APP_NAME="myapp"
VERSION=$(git describe --tags --always --dirty)
BUILD_TIME=$(date '+%Y-%m-%d_%H:%M:%S')

# 构建参数
LD_FLAGS="-w -s -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}"
TAGS="jsoniter"

# 执行构建
echo "Building ${APP_NAME}..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -buildvcs=false \
    -ldflags "${LD_FLAGS}" \
    -tags "${TAGS}" \
    -o "${APP_NAME}" \
    ./cmd/app # main.go文件目录

# 使用UPX压缩
upx --lzma --best "${APP_NAME}"