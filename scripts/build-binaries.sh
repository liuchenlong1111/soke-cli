#!/bin/bash

# 多平台编译脚本
# 用于编译 soke-cli 的所有平台二进制文件

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 获取版本号
VERSION=$(node -p "require('./package.json').version")
echo -e "${GREEN}开始编译 soke-cli v${VERSION}${NC}"

# 编译参数
LDFLAGS="-s -w -X main.Version=${VERSION}"

# 创建输出目录
BIN_DIR="bin"
mkdir -p ${BIN_DIR}

# 清理旧文件
echo -e "${YELLOW}清理旧的编译文件...${NC}"
rm -f ${BIN_DIR}/soke-cli-*

# 编译函数
build() {
    local GOOS=$1
    local GOARCH=$2
    local OUTPUT=$3

    echo -e "${YELLOW}编译 ${GOOS}-${GOARCH}...${NC}"

    GOOS=${GOOS} GOARCH=${GOARCH} go build \
        -ldflags "${LDFLAGS}" \
        -o ${BIN_DIR}/${OUTPUT} \
        .

    if [ $? -eq 0 ]; then
        local SIZE=$(du -h ${BIN_DIR}/${OUTPUT} | cut -f1)
        echo -e "${GREEN}✓ ${OUTPUT} (${SIZE})${NC}"
    else
        echo -e "${RED}✗ ${OUTPUT} 编译失败${NC}"
        exit 1
    fi
}

# 编译所有平台
echo ""
echo "编译目标平台:"

# macOS Intel
build "darwin" "amd64" "soke-cli-darwin-amd64"

# macOS Apple Silicon
build "darwin" "arm64" "soke-cli-darwin-arm64"

# Linux x64
build "linux" "amd64" "soke-cli-linux-amd64"

# Linux ARM64 (可选)
# build "linux" "arm64" "soke-cli-linux-arm64"

# Windows x64
build "windows" "amd64" "soke-cli-windows-amd64.exe"

echo ""
echo -e "${GREEN}所有平台编译完成!${NC}"
echo ""
echo "编译产物:"
ls -lh ${BIN_DIR}/soke-cli-*

echo ""
echo -e "${YELLOW}下一步:${NC}"
echo "1. 创建 Git 标签: git tag v${VERSION}"
echo "2. 推送标签: git push origin v${VERSION}"
echo "3. 上传到 GitHub Releases:"
echo "   gh release create v${VERSION} \\"
echo "     ${BIN_DIR}/soke-cli-darwin-amd64 \\"
echo "     ${BIN_DIR}/soke-cli-darwin-arm64 \\"
echo "     ${BIN_DIR}/soke-cli-linux-amd64 \\"
echo "     ${BIN_DIR}/soke-cli-windows-amd64.exe \\"
echo "     --title \"v${VERSION}\" \\"
echo "     --notes \"Release v${VERSION}\""
echo "4. 发布到 NPM: npm publish --access public"
