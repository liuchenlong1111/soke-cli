# 测试指南

## 快速开始

### 运行所有测试
```bash
make test-all
```

### 分模块测试

#### 1. 通讯录模块
```bash
make test-contact
```
测试内容：
- 获取部门列表
- 获取部门详情
- 获取部门用户列表
- 获取讲师列表
- 获取用户组列表
- 搜索用户

#### 2. 课程模块
```bash
make test-course
```
测试内容：
- 获取课程列表
- 获取课程分类
- 获取课程章节

#### 3. 考试模块
```bash
make test-exam
```
测试内容：
- 获取考试列表
- 获取考试分类

#### 4. 证书模块
```bash
make test-certificate
```
测试内容：
- 获取证书列表
- 获取证书分类

#### 5. 学分模块
```bash
make test-credit
```
测试内容：
- 获取学分日志

#### 6. 积分模块
```bash
make test-point
```
测试内容：
- 获取积分日志

#### 7. 线下培训模块
```bash
make test-training
```
测试内容：
- 获取培训列表
- 获取培训分类

#### 8. 学习地图模块
```bash
make test-learning-map
```
测试内容：
- 获取学习地图列表
- 获取学习地图分类

#### 9. 新闻模块
```bash
make test-news
```
测试内容：
- 获取新闻列表
- 获取新闻分类

#### 10. 作业模块
```bash
make test-clock
```
测试内容：
- 获取作业列表

#### 11. 素材库模块
```bash
make test-file
```
测试内容：
- 获取素材库列表
- 获取素材库分类

## 测试统计

### 总体情况
- **总CLI命令数**: 61个
- **E2E测试覆盖**: 24个 (39.3%)
- **测试模块数**: 11个

### 各模块测试数量
| 模块 | CLI总数 | 测试数量 | 覆盖率 |
|------|---------|----------|--------|
| contact | 15 | 6 | 40% |
| course | 8 | 3 | 37.5% |
| exam | 4 | 2 | 50% |
| certificate | 3 | 2 | 66.7% |
| credit | 1 | 1 | 100% |
| point | 4 | 1 | 25% |
| training | 6 | 2 | 33.3% |
| learning-map | 10 | 2 | 20% |
| news | 3 | 2 | 66.7% |
| clock | 3 | 1 | 33.3% |
| file | 4 | 2 | 50% |

## 测试策略

### 测试原则
1. ✓ 测试所有不需要参数的list接口
2. ✓ 测试所有categories接口
3. ✓ 测试可以用默认参数的get接口
4. ✗ 不测试需要动态ID的接口
5. ✗ 不测试写入/修改/删除操作
6. ✗ 不测试可能影响生产数据的接口

### 为什么不测试所有命令？
E2E测试主要验证核心功能可用性，而不是完整的功能测试。原因如下：

1. **需要动态ID的接口** (15个)
   - 需要先获取真实ID才能测试
   - 测试复杂度高，维护成本大
   - 示例：`contact +get-user --dept-user-id xxx`

2. **写入/修改/删除操作** (10个)
   - 可能影响生产数据
   - 需要清理测试数据
   - 示例：`contact position-create`, `training +create`

3. **需要特定参数的详细接口** (12个)
   - 需要复杂的测试数据准备
   - 测试结果不稳定
   - 示例：`course +list-lesson-faces --lesson-id xxx`

## 手动测试

对于未覆盖的命令，建议手动测试：

### 1. 获取详情类命令
```bash
# 先获取列表
./soke-cli course +list-courses --start-time 1672502400000 --end-time 1704038400000

# 使用返回的course-id测试详情
./soke-cli course +get-course --course-id <course-id>
```

### 2. 创建/修改类命令
```bash
# 创建职位
./soke-cli contact position-create --name "测试职位" --description "测试描述"

# 更新职位
./soke-cli contact position-update --position-id <id> --name "新名称"

# 删除职位
./soke-cli contact position-delete --position-id <id>
```

### 3. 积分操作
```bash
# 查询个人积分
./soke-cli point +get-user-info --dept-user-id <user-id>

# 添加积分
./soke-cli point +update-consume --trade-no "TEST001" --dept-user-id <user-id> --title "测试" --point 100

# 减少积分
./soke-cli point +update-consume --trade-no "TEST002" --dept-user-id <user-id> --title "测试" --point -50
```

## 故障排查

### 测试失败
如果测试失败，检查以下几点：

1. **登录状态**
   ```bash
   ./soke-cli config show
   ```

2. **网络连接**
   ```bash
   curl -I https://oapi.soke.cn
   ```

3. **查看详细错误**
   ```bash
   # 直接运行命令查看错误信息
   ./soke-cli <module> <command> <args>
   ```

### 常见问题

**Q: 测试提示未登录？**
```bash
./soke-cli config init
./soke-cli auth login
```

**Q: 某个模块测试失败？**
```bash
# 单独运行该模块测试查看详细信息
make test-<module>

# 或直接运行命令
./soke-cli <module> <command> <args>
```

**Q: 如何添加新的测试？**

编辑 `scripts/e2e-test.sh`，在对应模块添加测试用例：
```bash
test_command "module" "command" "args" "description"
```

## 持续集成

在CI/CD流程中使用：

```yaml
# GitHub Actions 示例
- name: Run tests
  run: |
    make test-all
```

```yaml
# GitLab CI 示例
test:
  script:
    - make test-all
```

## 发版前检查清单

- [ ] 运行 `make test` (单元测试)
- [ ] 运行 `make test-all` (E2E测试)
- [ ] 手动测试新增功能
- [ ] 检查 Git 状态干净
- [ ] 更新版本号
- [ ] 运行 `./scripts/release.sh`
