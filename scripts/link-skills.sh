#!/bin/bash

# soke-cli Skills 本地测试脚本
# 将本地 skills 链接到 Claude 的 skills 目录进行测试

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  soke-cli Skills 本地测试${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 检查 skills 目录是否存在
if [ ! -d "./skills" ]; then
    echo -e "${RED}错误: 未找到 ./skills 目录${NC}"
    exit 1
fi

# 查找所有可能的 skills 目录
echo -e "${YELLOW}查找 Claude skills 目录...${NC}"
SKILLS_DIRS=(
    "$HOME/.codex/skills"
    "$HOME/.openclaw/skills"
    "$HOME/.agents/skills"
    "$HOME/.workclaw/skills"
    "$HOME/.skills"
)

FOUND_DIRS=()
for dir in "${SKILLS_DIRS[@]}"; do
    if [ -d "$dir" ]; then
        FOUND_DIRS+=("$dir")
        echo -e "  ${GREEN}✓${NC} 找到: $dir"
    fi
done

if [ ${#FOUND_DIRS[@]} -eq 0 ]; then
    echo -e "${RED}错误: 未找到任何 Claude skills 目录${NC}"
    echo -e "${YELLOW}提示: 请先运行 'npx skills add' 安装任意一个 skill 来初始化目录${NC}"
    exit 1
fi
echo ""

# 选择要链接的目录
if [ ${#FOUND_DIRS[@]} -eq 1 ]; then
    TARGET_DIR="${FOUND_DIRS[0]}"
    echo -e "${BLUE}将链接到: ${TARGET_DIR}${NC}"
else
    echo -e "${YELLOW}找到多个 skills 目录，请选择:${NC}"
    for i in "${!FOUND_DIRS[@]}"; do
        echo -e "  $((i+1)). ${FOUND_DIRS[$i]}"
    done
    echo ""

    # 如果有参数，使用参数作为选择
    if [ -n "$1" ]; then
        choice=$1
    else
        echo -n "请输入序号 (1-${#FOUND_DIRS[@]}) 或 'all' 链接到所有目录: "
        read -r choice
    fi

    if [ "$choice" = "all" ]; then
        echo -e "${BLUE}将链接到所有目录${NC}"
        LINK_ALL=true
    elif [[ "$choice" =~ ^[0-9]+$ ]] && [ "$choice" -ge 1 ] && [ "$choice" -le ${#FOUND_DIRS[@]} ]; then
        TARGET_DIR="${FOUND_DIRS[$((choice-1))]}"
        echo -e "${BLUE}已选择: ${TARGET_DIR}${NC}"
        LINK_ALL=false
    else
        echo -e "${RED}无效的选择${NC}"
        exit 1
    fi
fi
echo ""

# 获取当前项目的绝对路径
PROJECT_DIR=$(pwd)
SOURCE_SKILLS_DIR="${PROJECT_DIR}/skills"

# 列出要链接的 skills
echo -e "${YELLOW}准备链接以下 skills:${NC}"
for skill_dir in "$SOURCE_SKILLS_DIR"/*; do
    if [ -d "$skill_dir" ]; then
        skill_name=$(basename "$skill_dir")
        echo -e "  - ${BLUE}${skill_name}${NC}"
    fi
done
echo ""

# 确认操作
echo -e "${YELLOW}是否继续? (y/n)${NC}"
read -r CONFIRM
if [ "$CONFIRM" != "y" ] && [ "$CONFIRM" != "Y" ]; then
    echo -e "${YELLOW}取消操作${NC}"
    exit 0
fi
echo ""

# 创建符号链接
echo -e "${YELLOW}创建符号链接...${NC}"
SUCCESS_COUNT=0
SKIP_COUNT=0

# 如果选择链接到所有目录
if [ "$LINK_ALL" = true ]; then
    for target_dir in "${FOUND_DIRS[@]}"; do
        echo -e "${BLUE}链接到: ${target_dir}${NC}"
        for skill_dir in "$SOURCE_SKILLS_DIR"/*; do
            if [ -d "$skill_dir" ]; then
                skill_name=$(basename "$skill_dir")
                target_link="${target_dir}/${skill_name}"

                if [ -L "$target_link" ]; then
                    current_target=$(readlink "$target_link")
                    if [ "$current_target" = "$skill_dir" ]; then
                        echo -e "  ${GREEN}✓${NC} ${skill_name} (已存在)"
                        SKIP_COUNT=$((SKIP_COUNT + 1))
                    else
                        rm "$target_link"
                        ln -s "$skill_dir" "$target_link"
                        echo -e "  ${GREEN}✓${NC} ${skill_name} (已更新)"
                        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
                    fi
                elif [ -e "$target_link" ]; then
                    echo -e "  ${RED}✗${NC} ${skill_name} (存在同名文件)"
                else
                    ln -s "$skill_dir" "$target_link"
                    echo -e "  ${GREEN}✓${NC} ${skill_name} (已链接)"
                    SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
                fi
            fi
        done
        echo ""
    done
else
    # 链接到单个目录
    for skill_dir in "$SOURCE_SKILLS_DIR"/*; do
        if [ -d "$skill_dir" ]; then
            skill_name=$(basename "$skill_dir")
            target_link="${TARGET_DIR}/${skill_name}"

            # 检查是否已存在
            if [ -L "$target_link" ]; then
                # 已存在符号链接，检查是否指向正确位置
                current_target=$(readlink "$target_link")
                if [ "$current_target" = "$skill_dir" ]; then
                    echo -e "  ${GREEN}✓${NC} ${skill_name} (已存在，指向正确)"
                    SKIP_COUNT=$((SKIP_COUNT + 1))
                else
                    echo -e "  ${YELLOW}!${NC} ${skill_name} (已存在，但指向: ${current_target})"
                    echo -n "    是否覆盖? (y/n): "
                    read -r overwrite
                    if [ "$overwrite" = "y" ] || [ "$overwrite" = "Y" ]; then
                        rm "$target_link"
                        ln -s "$skill_dir" "$target_link"
                        echo -e "    ${GREEN}✓${NC} 已覆盖"
                        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
                    else
                        echo -e "    ${YELLOW}跳过${NC}"
                        SKIP_COUNT=$((SKIP_COUNT + 1))
                    fi
                fi
            elif [ -e "$target_link" ]; then
                # 存在同名文件/目录
                echo -e "  ${RED}✗${NC} ${skill_name} (存在同名文件/目录)"
                echo -e "    请手动删除: rm -rf ${target_link}"
            else
                # 创建新链接
                ln -s "$skill_dir" "$target_link"
                echo -e "  ${GREEN}✓${NC} ${skill_name} (已链接)"
                SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
            fi
        fi
    done
fi
echo ""

# 完成
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  Skills 链接完成! ✓${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${BLUE}统计:${NC}"
echo -e "  新建链接: ${GREEN}${SUCCESS_COUNT}${NC}"
echo -e "  已存在: ${YELLOW}${SKIP_COUNT}${NC}"
echo ""
echo -e "${BLUE}下一步:${NC}"
echo -e "  1. 确保全局 CLI 已更新: ${BLUE}bash ./scripts/local-test.sh${NC}"
echo -e "  2. 在 Claude Code 中测试:"
echo -e "     ${BLUE}\"查询张三的学习档案\"${NC}"
echo -e "     ${BLUE}\"查询技术部的学员学习情况\"${NC}"
echo ""
echo -e "${YELLOW}提示:${NC}"
echo -e "  - Skills 已链接到: ${TARGET_DIR}"
echo -e "  - 修改本地 skills 文件会立即生效"
echo -e "  - 删除链接: ${BLUE}rm ${TARGET_DIR}/soke-*${NC}"
echo ""
