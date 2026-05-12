# soke-cli 测试指南

## 测试类型

### 1. 单元测试（Unit Tests）
测试单个函数和模块的逻辑正确性，不依赖外部 API。

### 2. 端到端测试（E2E Tests）
测试完整的命令执行流程，会真实调用 API 接口。

---

## 快速开始

### 运行所有测试
```bash
make test-all
```

### 只运行单元测试
```bash
make test
# 或
go test -v ./...
```

### 只运行端到端测试
```bash
make e2e-test
# 或
./scripts/e2e-test.sh
```

---

## 端到端测试说明

### 前置条件
1. **已配置应用凭证**
   ```bash
   ./soke-cli config init
   ```

2. **已完成用户登录**
   ```bash
   ./soke-cli auth login
   ```

3. **已编译二进制文件**
   ```bash
   go build -o soke-cli main.go
   ```

### 测试覆盖的模块
- ✅ Contact（通讯录）
- ✅ Course（课程）
- ✅ Exam（考试）
- ✅ Certificate（证书）
- ✅ Credit（学分）
- ✅ Point（积分）
- ✅ Training（培训）
- ✅ Learning Map（学习地图）
- ✅ News（新闻公告）

### 测试输出示例
```
========================================
  soke-cli E2E 测试 v1.0.1
========================================

✓ 已登录

开始测试...

[Contact 模块]
测试 [contact +list-departments]:  获取部门列表 ... ✓ 通过
测试 [contact +get-department]:  获取部门详情 ... ✓ 通过
测试 [contact +list-department-users]:  获取部门用户列表 ... ✓ 通过

[Course 模块]
测试 [course +list-courses]:  获取课程列表 ... ✓ 通过
测试 [course +list-categories]:  获取课程分类 ... ✓ 通过

========================================
  测试结果汇总
========================================

总测试数: 25
通过: 25
失败: 0

========================================
  所有测试通过! ✓
========================================
```

---

## 发布前测试流程

### 方法一：使用发布脚本（推荐）
```bash
./scripts/release.sh
```
发布脚本会自动运行所有测试，测试失败会阻止发布。

### 方法二：手动测试
```bash
# 1. 运行单元测试
make test

# 2. 运行端到端测试
make e2e-test

# 3. 如果都通过，继续发布流程
./scripts/build-binaries.sh
git tag v1.0.x
# ...
```

---

## 添加新的测试

### 添加单元测试
在对应模块目录下创建 `*_test.go` 文件：

```go
// shortcuts/contact/contact_list_users_test.go
package contact

import (
    "context"
    "testing"
)

func TestContactListUsers_DryRun(t *testing.T) {
    // 测试代码
}
```

### 添加端到端测试
编辑 `scripts/e2e-test.sh`，添加新的测试命令：

```bash
test_command "module" "+command" "--arg value" "描述"
```

---

## 常见问题

### Q: 端到端测试失败怎么办？
A: 检查以下几点：
1. 是否已登录：`./soke-cli config show`
2. Token 是否过期：重新运行 `./soke-cli auth login`
3. API 地址是否正确：检查配置文件
4. 网络是否正常：尝试手动执行失败的命令

### Q: 如何跳过某些测试？
A: 在 `scripts/e2e-test.sh` 中注释掉对应的 `test_command` 行。

### Q: 测试数据从哪里来？
A: 端到端测试使用真实的生产/测试环境数据，确保测试账号有足够的权限和数据。

---

## 持续集成（CI）

如果使用 GitHub Actions 或其他 CI 工具，可以添加以下配置：

```yaml
# .github/workflows/test.yml
name: Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: '1.25'
      - name: Run unit tests
        run: make test
      # 注意：E2E 测试需要配置凭证，不适合在公开 CI 中运行
```

---

## 最佳实践

1. **每次修改代码后都运行测试**
   ```bash
   make test-all
   ```

2. **发布前必须运行完整测试**
   ```bash
   ./scripts/release.sh  # 自动包含测试
   ```

3. **添加新功能时同步添加测试**
   - 新增命令 → 添加单元测试
   - 新增模块 → 添加端到端测试

4. **定期更新测试数据**
   - 确保测试使用的 ID、时间戳等参数是有效的
