#!/bin/bash

# soke-cli 端到端测试脚本
# 测试所有命令是否能正常执行
# 支持按模块测试: ./e2e-test.sh [module_name]

# 注意：不使用 set -e，让所有测试都能运行完

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 测试结果统计
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0
FAILED_COMMANDS=()

# 获取版本号
VERSION=$(node -p "require('./package.json').version")

# 获取要测试的模块（如果指定）
TEST_MODULE=$1

echo -e "${BLUE}========================================${NC}"
if [ -z "$TEST_MODULE" ]; then
    echo -e "${BLUE}  soke-cli E2E 测试 v${VERSION}${NC}"
else
    echo -e "${BLUE}  soke-cli E2E 测试 - ${TEST_MODULE} 模块${NC}"
fi
echo -e "${BLUE}========================================${NC}"
echo ""

# 检查是否已编译
if [ ! -f "./soke-cli" ]; then
    echo -e "${YELLOW}未找到编译文件，开始编译...${NC}"
    go build -o soke-cli main.go
    echo -e "${GREEN}✓ 编译完成${NC}"
    echo ""
fi

# 检查全局安装的 CLI 版本
echo -e "${YELLOW}检查全局 CLI 安装...${NC}"
GLOBAL_CLI=$(which soke-cli 2>/dev/null)
if [ -n "$GLOBAL_CLI" ]; then
    echo -e "  全局 CLI 路径: ${BLUE}${GLOBAL_CLI}${NC}"

    # 检查全局版本是否支持 learning-profile
    if ! $GLOBAL_CLI learning-profile --help &>/dev/null; then
        echo -e "${YELLOW}  警告: 全局 CLI 不支持 learning-profile 模块${NC}"
        echo -e "${YELLOW}  需要更新全局安装以支持 skill 测试${NC}"
        echo ""
        echo -e "${YELLOW}是否要将本地编译版本安装到全局? (y/n)${NC}"
        read -r INSTALL_GLOBAL
        if [ "$INSTALL_GLOBAL" = "y" ] || [ "$INSTALL_GLOBAL" = "Y" ]; then
            echo -e "${YELLOW}安装到全局...${NC}"
            sudo cp ./soke-cli $GLOBAL_CLI
            echo -e "${GREEN}✓ 全局 CLI 已更新${NC}"
        else
            echo -e "${YELLOW}跳过全局安装，将仅测试本地版本${NC}"
        fi
    else
        echo -e "${GREEN}✓ 全局 CLI 支持 learning-profile 模块${NC}"
    fi
else
    echo -e "${YELLOW}  未找到全局 CLI 安装${NC}"
    echo -e "${YELLOW}  Skill 测试需要全局安装 soke-cli${NC}"
fi
echo ""

# 检查是否已登录
echo -e "${YELLOW}检查登录状态...${NC}"
if ! ./soke-cli config show &>/dev/null; then
    echo -e "${RED}错误: 未配置或未登录${NC}"
    echo "请先运行: ./soke-cli config init && ./soke-cli auth login"
    exit 1
fi
echo -e "${GREEN}✓ 已登录${NC}"
echo ""

# 测试函数
test_command() {
    local module=$1
    local command=$2
    local args=$3
    local description=$4

    TOTAL_TESTS=$((TOTAL_TESTS + 1))

    echo -n "测试 [$module $command]: $description ... "

    if ./soke-cli $module $command $args &>/dev/null; then
        echo -e "${GREEN}✓ 通过${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        return 0
    else
        echo -e "${RED}✗ 失败${NC}"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        FAILED_COMMANDS+=("$module $command $args")
        return 1
    fi
}

# 判断是否应该运行某个模块的测试
should_test_module() {
    local module=$1
    if [ -z "$TEST_MODULE" ]; then
        return 0  # 没有指定模块，测试所有
    elif [ "$TEST_MODULE" = "$module" ]; then
        return 0  # 指定的模块匹配
    else
        return 1  # 不匹配，跳过
    fi
}

echo -e "${BLUE}开始测试...${NC}"
echo ""

# ==================== Contact 模块 ====================
if should_test_module "contact"; then
    echo -e "${YELLOW}[Contact 模块]${NC}"
    test_command "contact" "+list-departments" "" "获取部门列表"
    test_command "contact" "+get-department" "--dept-id 1" "获取部门详情"
    test_command "contact" "+list-department-users" "--dept-id 1" "获取部门用户列表"
    test_command "contact" "+list-lectors" "" "获取讲师列表"
    test_command "contact" "+list-groups" "--start-time 1672502400000 --end-time 1704038400000" "获取用户组列表"
    test_command "contact" "+search-user" "--dept-user-name 测试" "搜索用户"
    test_command "contact" "+search-dept" "--dept-name 测试" "搜索部门"
    echo ""
fi

# ==================== Learning Profile 模块 ====================
if should_test_module "learning-profile"; then
    echo -e "${YELLOW}[Learning Profile 模块]${NC}"
    test_command "learning-profile" "+list" "--offset 0 --page-size 10" "获取学习档案列表"
    test_command "learning-profile" "+list" "--is-new 1 --page-size 5" "获取新员工学习档案"
    echo ""
fi

# ==================== Course 模块 ====================
if should_test_module "course"; then
    echo -e "${YELLOW}[Course 模块]${NC}"
    test_command "course" "+list-courses" "--start-time 1672502400000 --end-time 1704038400000" "获取课程列表"
    test_command "course" "+list-categories" "" "获取课程分类"
    test_command "course" "+list-lessons" "--course-id F4025496-6566-4ACD-8CD4-06E3D55BBA79" "获取课程章节"
    echo ""
fi

# ==================== Exam 模块 ====================
if should_test_module "exam"; then
    echo -e "${YELLOW}[Exam 模块]${NC}"
    test_command "exam" "+list-exams" "--start-time 1672502400000 --end-time 1704038400000" "获取考试列表"
    test_command "exam" "+list-categories" "" "获取考试分类"
    echo ""
fi

# ==================== Certificate 模块 ====================
if should_test_module "certificate"; then
    echo -e "${YELLOW}[Certificate 模块]${NC}"
    test_command "certificate" "+list-certificates" "--start-time 1672502400000 --end-time 1704038400000" "获取证书列表"
    test_command "certificate" "+list-categories" "" "获取证书分类"
    echo ""
fi

# ==================== Credit 模块 ====================
if should_test_module "credit"; then
    echo -e "${YELLOW}[Credit 模块]${NC}"
    test_command "credit" "+list-logs" "--start-time 1672502400000 --end-time 1704038400000" "获取学分日志"
    echo ""
fi

# ==================== Point 模块 ====================
if should_test_module "point"; then
    echo -e "${YELLOW}[Point 模块]${NC}"
    test_command "point" "+list-logs" "--start-time 1672502400000 --end-time 1704038400000" "获取积分日志"
    echo ""
fi

# ==================== Training 模块 ====================
if should_test_module "training"; then
    echo -e "${YELLOW}[Training 模块]${NC}"
    test_command "training" "+list-trainings" "--start-time 1672502400000 --end-time 1704038400000" "获取培训列表"
    test_command "training" "+list-categories" "" "获取培训分类"
    echo ""
fi

# ==================== Learning Map 模块 ====================
if should_test_module "learning-map"; then
    echo -e "${YELLOW}[Learning Map 模块]${NC}"
    test_command "learning-map" "+list-maps" "--start-time 1672502400000 --end-time 1704038400000" "获取学习地图列表"
    test_command "learning-map" "+list-categories" "" "获取学习地图分类"
    echo ""
fi

# ==================== News 模块 ====================
if should_test_module "news"; then
    echo -e "${YELLOW}[News 模块]${NC}"
    test_command "news" "+list-news" "" "获取新闻列表"
    test_command "news" "+list-categories" "" "获取新闻分类"
    echo ""
fi

# ==================== Clock 模块 ====================
if should_test_module "clock"; then
    echo -e "${YELLOW}[Clock 模块]${NC}"
    test_command "clock" "list-learnings" "--page 1 --page-size 10" "获取作业列表"
    echo ""
fi

# ==================== File 模块 ====================
if should_test_module "file"; then
    echo -e "${YELLOW}[File 模块]${NC}"
    test_command "file" "list-files" "" "获取素材库列表"
    test_command "file" "list-categories" "" "获取素材库分类"
    echo ""
fi

# ==================== 测试结果汇总 ====================
echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  测试结果汇总${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "总测试数: ${BLUE}${TOTAL_TESTS}${NC}"
echo -e "通过: ${GREEN}${PASSED_TESTS}${NC}"
echo -e "失败: ${RED}${FAILED_TESTS}${NC}"
echo ""

if [ ${FAILED_TESTS} -gt 0 ]; then
    echo -e "${RED}失败的命令:${NC}"
    for cmd in "${FAILED_COMMANDS[@]}"; do
        echo -e "  ${RED}✗${NC} ./soke-cli $cmd"
    done
    echo ""
    echo -e "${RED}测试失败! 请修复以上问题后再发布。${NC}"
    exit 1
else
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}  所有测试通过! ✓${NC}"
    echo -e "${GREEN}========================================${NC}"
    exit 0
fi
