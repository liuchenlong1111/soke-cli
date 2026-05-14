#!/bin/bash

# soke-cli 本地测试脚本
# 用于本地开发测试：编译 -> 安装到全局 -> 验证功能

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  soke-cli 本地测试${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 步骤1: 编译
echo -e "${YELLOW}[1/4] 编译 soke-cli...${NC}"
if go build -o soke-cli main.go; then
    echo -e "${GREEN}✓ 编译成功${NC}"
else
    echo -e "${RED}✗ 编译失败${NC}"
    exit 1
fi
echo ""

# 步骤2: 检查本地版本功能
echo -e "${YELLOW}[2/4] 检查本地版本功能...${NC}"
echo -e "  检查 learning-profile 模块..."
if ./soke-cli learning-profile --help &>/dev/null; then
    echo -e "${GREEN}  ✓ learning-profile 模块存在${NC}"
else
    echo -e "${RED}  ✗ learning-profile 模块不存在${NC}"
    exit 1
fi

echo -e "  检查 contact +search-dept 命令..."
if ./soke-cli contact +search-dept --help &>/dev/null; then
    echo -e "${GREEN}  ✓ contact +search-dept 命令存在${NC}"
else
    echo -e "${RED}  ✗ contact +search-dept 命令不存在${NC}"
    exit 1
fi
echo ""

# 步骤3: 安装到全局
echo -e "${YELLOW}[3/4] 安装到全局...${NC}"

# 查找全局 CLI 路径
GLOBAL_CLI=$(which soke-cli 2>/dev/null)
if [ -z "$GLOBAL_CLI" ]; then
    echo -e "${YELLOW}  未找到全局 soke-cli 安装${NC}"
    echo -e "${YELLOW}  请先通过 npm 安装: npm install -g @sokeai/cli${NC}"
    echo -e "${YELLOW}  跳过全局安装，仅测试本地版本${NC}"
    SKIP_GLOBAL=true
else
    echo -e "  全局 CLI 路径: ${BLUE}${GLOBAL_CLI}${NC}"

    # 检查是否需要更新
    NEED_UPDATE=false
    if ! $GLOBAL_CLI learning-profile --help &>/dev/null; then
        echo -e "${YELLOW}  全局版本不支持 learning-profile 模块${NC}"
        NEED_UPDATE=true
    elif ! $GLOBAL_CLI contact +search-dept --help &>/dev/null; then
        echo -e "${YELLOW}  全局版本不支持 contact +search-dept 命令${NC}"
        NEED_UPDATE=true
    fi

    if [ "$NEED_UPDATE" = true ]; then
        echo -e "${YELLOW}  需要更新全局安装${NC}"
        echo ""
        echo -e "${YELLOW}  是否要用本地版本覆盖全局安装? (y/n)${NC}"
        read -r CONFIRM

        if [ "$CONFIRM" = "y" ] || [ "$CONFIRM" = "Y" ]; then
            # 备份原文件
            BACKUP_FILE="${GLOBAL_CLI}.backup.$(date +%Y%m%d_%H%M%S)"
            echo -e "  备份原文件到: ${BACKUP_FILE}"
            if sudo cp $GLOBAL_CLI $BACKUP_FILE; then
                echo -e "${GREEN}  ✓ 备份成功${NC}"
            else
                echo -e "${RED}  ✗ 备份失败${NC}"
                exit 1
            fi

            # 安装新版本
            echo -e "  安装新版本..."
            if sudo cp ./soke-cli $GLOBAL_CLI; then
                echo -e "${GREEN}  ✓ 安装成功${NC}"
                echo -e "${BLUE}  备份文件: ${BACKUP_FILE}${NC}"
            else
                echo -e "${RED}  ✗ 安装失败${NC}"
                echo -e "${YELLOW}  恢复备份...${NC}"
                sudo cp $BACKUP_FILE $GLOBAL_CLI
                exit 1
            fi
            SKIP_GLOBAL=false
        else
            echo -e "${YELLOW}  跳过全局安装${NC}"
            SKIP_GLOBAL=true
        fi
    else
        echo -e "${GREEN}  ✓ 全局版本已是最新，无需更新${NC}"
        SKIP_GLOBAL=false
    fi
fi
echo ""

# 步骤4: 运行测试
echo -e "${YELLOW}[4/4] 运行功能测试...${NC}"

# 测试本地版本
echo -e "  ${BLUE}测试本地版本:${NC}"
if ./soke-cli learning-profile +list --help &>/dev/null; then
    echo -e "${GREEN}  ✓ learning-profile +list 命令可用${NC}"
else
    echo -e "${RED}  ✗ learning-profile +list 命令不可用${NC}"
    exit 1
fi

if ./soke-cli contact +search-dept --help &>/dev/null; then
    echo -e "${GREEN}  ✓ contact +search-dept 命令可用${NC}"
else
    echo -e "${RED}  ✗ contact +search-dept 命令不可用${NC}"
    exit 1
fi

if ./soke-cli contact +search-user --help &>/dev/null; then
    echo -e "${GREEN}  ✓ contact +search-user 命令可用${NC}"
else
    echo -e "${RED}  ✗ contact +search-user 命令不可用${NC}"
    exit 1
fi

# 测试全局版本（如果已安装）
if [ "$SKIP_GLOBAL" = false ] && [ -n "$GLOBAL_CLI" ]; then
    echo ""
    echo -e "  ${BLUE}测试全局版本:${NC}"
    if $GLOBAL_CLI learning-profile +list --help &>/dev/null; then
        echo -e "${GREEN}  ✓ learning-profile +list 命令可用${NC}"
    else
        echo -e "${RED}  ✗ learning-profile +list 命令不可用${NC}"
        echo -e "${YELLOW}  提示: 全局安装可能未成功更新${NC}"
    fi

    if $GLOBAL_CLI contact +search-dept --help &>/dev/null; then
        echo -e "${GREEN}  ✓ contact +search-dept 命令可用${NC}"
    else
        echo -e "${RED}  ✗ contact +search-dept 命令不可用${NC}"
        echo -e "${YELLOW}  提示: 全局安装可能未成功更新${NC}"
    fi

    if $GLOBAL_CLI contact +search-user --help &>/dev/null; then
        echo -e "${GREEN}  ✓ contact +search-user 命令可用${NC}"
    else
        echo -e "${RED}  ✗ contact +search-user 命令不可用${NC}"
        echo -e "${YELLOW}  提示: 全局安装可能未成功更新${NC}"
    fi
fi
echo ""

# 完成
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  本地测试完成! ✓${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${BLUE}下一步:${NC}"
if [ "$SKIP_GLOBAL" = true ]; then
    echo -e "  ${YELLOW}提示: 未更新全局安装，Skills 可能无法使用新功能${NC}"
    echo -e "  ${YELLOW}如需测试 Skills，请重新运行此脚本并选择更新全局安装${NC}"
    echo ""
fi
echo -e "  1. 运行完整测试: ${BLUE}bash ./scripts/e2e-test.sh${NC}"
echo -e "  2. 测试 Skills: ${BLUE}npx skills add liuchenlong1111/soke-cli -y -g${NC}"
echo -e "  3. 在 AI Agent 中测试: ${BLUE}\"查询张三的学习档案\"${NC}"
echo ""

