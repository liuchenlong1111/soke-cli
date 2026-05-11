#!/bin/bash

# 授客CLI双远程仓库推送脚本
# 用途：将代码同时推送到 GitHub 和 Aliyun 远程仓库

set -e  # 遇到错误立即退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的信息
print_info() {
    echo -e "${GREEN}[信息]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[警告]${NC} $1"
}

print_error() {
    echo -e "${RED}[错误]${NC} $1"
}

print_highlight() {
    echo -e "${BLUE}[提示]${NC} $1"
}

# 获取 package.json 中的当前版本号
get_current_version() {
    if [[ ! -f "package.json" ]]; then
        print_error "未找到 package.json 文件"
        exit 1
    fi

    # 使用 grep 和 sed 提取版本号
    version=$(grep '"version"' package.json | sed -E 's/.*"version": "([^"]+)".*/\1/')
    echo "$version"
}

# 更新 package.json 中的版本号
update_package_json_version() {
    local new_version=$1

    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS 使用 BSD sed
        sed -i '' "s/\"version\": \".*\"/\"version\": \"$new_version\"/" package.json
    else
        # Linux 使用 GNU sed
        sed -i "s/\"version\": \".*\"/\"version\": \"$new_version\"/" package.json
    fi

    print_info "已更新 package.json 版本号为: ${GREEN}${new_version}${NC}"
}

# 更新 internal/version/checker.go 中的版本号
update_version_checker() {
    local new_version=$1
    local version_file="internal/version/checker.go"

    if [[ ! -f "$version_file" ]]; then
        print_warning "未找到 $version_file 文件，跳过更新"
        return
    fi

    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS 使用 BSD sed
        sed -i '' "s/var Version = \".*\"/var Version = \"$new_version\"/" "$version_file"
    else
        # Linux 使用 GNU sed
        sed -i "s/var Version = \".*\"/var Version = \"$new_version\"/" "$version_file"
    fi

    print_info "已更新 $version_file 版本号为: ${GREEN}${new_version}${NC}"
}

# 更新所有版本号
update_all_versions() {
    local new_version=$1

    update_package_json_version "$new_version"
    update_version_checker "$new_version"
}

# 验证版本号格式（简单的语义化版本检查）
validate_version() {
    local version=$1
    if [[ ! $version =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        print_error "版本号格式不正确，应为 x.y.z 格式（如 1.0.4）"
        return 1
    fi
    return 0
}

# 获取当前分支名
CURRENT_BRANCH=$(git branch --show-current)

# 显示当前分支并确认
echo ""
print_info "当前所在分支: ${GREEN}${CURRENT_BRANCH}${NC}"
echo ""
read -p "是否在当前分支提交代码？(y/n): " confirm

if [[ "$confirm" != "y" && "$confirm" != "Y" ]]; then
    print_warning "用户取消操作，脚本退出"
    exit 0
fi

# 询问是否作为新版本发布
echo ""
read -p "本次改动是否作为新版本发布？(y/n): " is_release

IS_RELEASE=false
NEW_VERSION=""

if [[ "$is_release" == "y" || "$is_release" == "Y" ]]; then
    IS_RELEASE=true

    # 获取当前版本号
    CURRENT_VERSION=$(get_current_version)
    print_highlight "当前版本号: ${BLUE}${CURRENT_VERSION}${NC}"

    # 提示输入新版本号
    while true; do
        echo ""
        read -p "请输入新版本号 (格式: x.y.z): " NEW_VERSION

        if [[ -z "$NEW_VERSION" ]]; then
            print_error "版本号不能为空"
            continue
        fi

        if validate_version "$NEW_VERSION"; then
            if [[ "$NEW_VERSION" == "$CURRENT_VERSION" ]]; then
                print_warning "新版本号与当前版本号相同，是否继续？(y/n)"
                read -p "> " continue_same
                if [[ "$continue_same" == "y" || "$continue_same" == "Y" ]]; then
                    break
                fi
            else
                break
            fi
        fi
    done

    # 更新所有文件中的版本号
    update_all_versions "$NEW_VERSION"

    # 将版本更新的改动加入暂存区
    git add package.json internal/version/checker.go
    print_info "已将版本更新加入暂存区"
fi

# 执行 git add .
print_info "正在暂存所有更改..."
git add .

# 检查是否有需要提交的更改
if git diff --cached --quiet; then
    print_warning "没有需要提交的更改"
    exit 0
fi

# 提示输入 commit 信息
echo ""
if [[ "$IS_RELEASE" == true ]]; then
    # 如果是版本发布，提供默认的 commit 信息
    DEFAULT_COMMIT="chore: release v${NEW_VERSION}"
    print_highlight "建议的 commit 信息: ${BLUE}${DEFAULT_COMMIT}${NC}"
    read -p "请输入 commit 信息 (直接回车使用建议信息): " commit_info

    if [[ -z "$commit_info" ]]; then
        commit_info="$DEFAULT_COMMIT"
    fi
else
    read -p "请输入 commit 信息: " commit_info

    if [[ -z "$commit_info" ]]; then
        print_error "commit 信息不能为空"
        exit 1
    fi
fi

# 执行 commit
print_info "正在提交代码..."
git commit -m "$commit_info"

# 保存当前分支名（用于后续合并）
FEATURE_BRANCH=$CURRENT_BRANCH

# 切换到 master 分支
print_info "切换到 master 分支..."
git checkout master

# 合并功能分支到 master
print_info "合并 ${FEATURE_BRANCH} 到 master..."
git merge $FEATURE_BRANCH

# 推送到 origin (Aliyun) 的 master 分支
print_info "推送到 Aliyun 远程仓库 (origin master)..."
git push origin master

# 推送到 github 的 master 分支
print_info "推送到 GitHub 远程仓库 (github master)..."
git push github master

# 切换到 main 分支
print_info "切换到 main 分支..."
git checkout main

# 合并 master 到 main
print_info "合并 master 到 main..."
git merge master

# 推送到 github 的 main 分支
print_info "推送到 GitHub 远程仓库 (github main)..."
git push github main

# 切换回原分支
print_info "切换回原分支 ${FEATURE_BRANCH}..."
git checkout $FEATURE_BRANCH

echo ""
print_info "✅ 代码推送完成！"
echo ""
print_info "推送摘要:"
echo "  - 功能分支: ${FEATURE_BRANCH}"
echo "  - Commit 信息: ${commit_info}"

if [[ "$IS_RELEASE" == true ]]; then
    echo "  - 🎉 版本发布: ${GREEN}v${NEW_VERSION}${NC}"
fi

echo "  - Aliyun (origin): master 分支已更新"
echo "  - GitHub (github): master 和 main 分支已更新"

if [[ "$IS_RELEASE" == true ]]; then
    echo ""
    print_highlight "📦 下一步操作："
    echo "  1. 编译二进制文件: make build"
    echo "  2. 发布到 NPM: npm publish"
    echo "  3. 创建 GitHub Release (可选)"
fi

echo ""
