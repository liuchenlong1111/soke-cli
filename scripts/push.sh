#!/bin/bash

# 授客CLI双远程仓库推送脚本
# 用途：将代码同时推送到 GitHub 和 Aliyun 远程仓库

set -e  # 遇到错误立即退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
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
read -p "请输入 commit 信息: " commit_info

if [[ -z "$commit_info" ]]; then
    print_error "commit 信息不能为空"
    exit 1
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
echo "  - Aliyun (origin): master 分支已更新"
echo "  - GitHub (github): master 和 main 分支已更新"
echo ""
