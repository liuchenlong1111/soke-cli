#!/bin/bash

# soke-cli 自动化发布脚本
# 完整流程: 编译 -> 创建标签 -> 上传 GitHub Releases -> 发布 NPM

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 获取版本号
VERSION=$(node -p "require('./package.json').version")

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  soke-cli 发布脚本 v${VERSION}${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 检查必要的工具
check_tools() {
    echo -e "${YELLOW}检查必要工具...${NC}"

    if ! command -v go &> /dev/null; then
        echo -e "${RED}错误: 未安装 Go${NC}"
        exit 1
    fi

    if ! command -v gh &> /dev/null; then
        echo -e "${RED}错误: 未安装 GitHub CLI (gh)${NC}"
        echo "安装方法: brew install gh"
        exit 1
    fi

    if ! command -v npm &> /dev/null; then
        echo -e "${RED}错误: 未安装 npm${NC}"
        exit 1
    fi

    echo -e "${GREEN}✓ 所有工具已就绪${NC}"
    echo ""
}

# 检查 Git 状态
check_git_status() {
    echo -e "${YELLOW}检查 Git 状态...${NC}"

    if [ -n "$(git status --porcelain)" ]; then
        echo -e "${RED}错误: 工作目录有未提交的更改${NC}"
        echo "请先提交或暂存所有更改"
        git status --short
        exit 1
    fi

    echo -e "${GREEN}✓ 工作目录干净${NC}"
    echo ""
}

# 检查标签是否已存在
check_tag() {
    echo -e "${YELLOW}检查标签 v${VERSION}...${NC}"

    if git rev-parse "v${VERSION}" >/dev/null 2>&1; then
        echo -e "${RED}错误: 标签 v${VERSION} 已存在${NC}"
        echo "请更新 package.json 中的版本号"
        exit 1
    fi

    echo -e "${GREEN}✓ 标签不存在，可以继续${NC}"
    echo ""
}

# 运行测试
run_tests() {
    echo -e "${YELLOW}运行单元测试...${NC}"

    if ! go test ./...; then
        echo -e "${RED}错误: 单元测试失败${NC}"
        exit 1
    fi

    echo -e "${GREEN}✓ 单元测试通过${NC}"
    echo ""

    echo -e "${YELLOW}运行端到端测试...${NC}"

    if ! ./scripts/e2e-test.sh; then
        echo -e "${RED}错误: 端到端测试失败${NC}"
        exit 1
    fi

    echo -e "${GREEN}✓ 端到端测试通过${NC}"
    echo ""
}

# 编译所有平台
build_binaries() {
    echo -e "${YELLOW}编译所有平台...${NC}"

    if ! ./scripts/build-binaries.sh; then
        echo -e "${RED}错误: 编译失败${NC}"
        exit 1
    fi

    echo -e "${GREEN}✓ 编译完成${NC}"
    echo ""
}

# 创建并推送标签
create_tag() {
    echo -e "${YELLOW}创建 Git 标签 v${VERSION}...${NC}"

    git tag -a "v${VERSION}" -m "Release v${VERSION}"
    git push github "v${VERSION}"

    echo -e "${GREEN}✓ 标签已创建并推送${NC}"
    echo ""
}

# 上传到 GitHub Releases
upload_github_release() {
    echo -e "${YELLOW}上传到 GitHub Releases...${NC}"

    # 生成 Release Notes（可以自定义）
    RELEASE_NOTES="Release v${VERSION}

## 安装方法

\`\`\`bash
npm install -g @sokeai/cli
\`\`\`

## 更新日志

请查看 commit 历史了解详细更改。
"

    gh release create "v${VERSION}" \
        bin/soke-cli-darwin-amd64 \
        bin/soke-cli-darwin-arm64 \
        bin/soke-cli-linux-amd64 \
        bin/soke-cli-windows-amd64.exe \
        --title "v${VERSION}" \
        --notes "${RELEASE_NOTES}"

    echo -e "${GREEN}✓ GitHub Release 创建成功${NC}"
    echo ""
}

# 发布到 NPM
publish_npm() {
    echo -e "${YELLOW}发布到 NPM...${NC}"

    # 检查是否已登录
    if ! npm whoami &> /dev/null; then
        echo -e "${RED}错误: 未登录 NPM${NC}"
        echo "请先运行: npm login"
        exit 1
    fi

    # 发布
    npm publish --access public

    echo -e "${GREEN}✓ NPM 包发布成功${NC}"
    echo ""
}

# 清理编译产物
cleanup() {
    echo -e "${YELLOW}清理编译产物...${NC}"
    rm -rf bin/
    echo -e "${GREEN}✓ 清理完成${NC}"
    echo ""
}

# 主流程
main() {
    echo -e "${BLUE}开始发布流程...${NC}"
    echo ""

    check_tools
    check_git_status
    check_tag
    run_tests
    build_binaries
    create_tag
    upload_github_release
    publish_npm
    cleanup

    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}  发布成功! 🎉${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo -e "版本: ${GREEN}v${VERSION}${NC}"
    echo -e "NPM: ${BLUE}https://www.npmjs.com/package/@sokeai/cli${NC}"
    echo -e "GitHub: ${BLUE}https://github.com/sokeai/soke-cli/releases/tag/v${VERSION}${NC}"
    echo ""
    echo "安装或升级到最新版本:"
    echo -e "${YELLOW}npm install -g @sokeai/cli@latest${NC}"
}

# 确认发布
echo -e "${YELLOW}即将发布 soke-cli v${VERSION}${NC}"
echo ""
echo "发布流程包括:"
echo "  1. 检查工具和 Git 状态"
echo "  2. 运行测试"
echo "  3. 编译所有平台二进制文件"
echo "  4. 创建 Git 标签并推送"
echo "  5. 上传到 GitHub Releases"
echo "  6. 发布到 NPM"
echo ""
read -p "确认继续? (y/N) " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "已取消发布"
    exit 0
fi

echo ""
main
