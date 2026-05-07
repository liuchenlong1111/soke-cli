# Contact Shortcuts

通讯录相关的快捷命令实现。

## 文件说明

### contact_get_department.go
获取部门详情的命令实现。

**功能特性：**
- 通过部门 ID 获取部门详细信息
- 支持 user 和 bot 两种认证方式
- 支持多种输出格式（json/table/pretty）
- 支持 dry-run 模式预览请求

**使用示例：**
```bash
# 获取部门信息
soke-cli contact +get-department --dept-id 123456

# 使用 JSON 格式输出
soke-cli contact +get-department --dept-id 123456 --format json

# 预览请求（不实际执行）
soke-cli contact +get-department --dept-id 123456 --dry-run
```

**API 端点：**
- `GET https://oapi.soke.cn/oa/department/info`

**参数：**
- `--dept-id`: 部门 ID（必填）

**权限范围：**
- User: `contact:department:readonly`
- Bot: `contact:department:readonly`

**输出字段：**
- `dept_id`: 部门 ID
- `dept_name`: 部门名称
- `parent_id`: 父部门 ID
- `company_id`: 公司 ID

## 测试

### contact_get_department_test.go
包含以下测试用例：

1. **TestContactGetDepartment_DryRun** - 测试 dry-run 功能
2. **TestContactGetDepartment_Metadata** - 测试元数据配置
   - Service 名称
   - Command 名称
   - Risk 级别
   - HasFormat 标志
3. **TestContactGetDepartment_Flags** - 测试命令行参数
   - dept-id 参数存在性
   - dept-id 必填验证
4. **TestContactGetDepartment_Scopes** - 测试权限范围
   - UserScopes 配置
   - BotScopes 配置
5. **TestContactGetDepartment_AuthTypes** - 测试认证类型
   - 支持 user 认证
   - 支持 bot 认证
6. **TestContactGetDepartment_Execute** - 测试执行函数存在性

### 运行测试

```bash
# 运行所有 contact 包测试
go test ./shortcuts/contact/... -v

# 运行特定测试
go test ./shortcuts/contact/contact_get_department_test.go -v

# 查看测试覆盖率
go test ./shortcuts/contact/... -cover
```

## 开发规范

1. 所有文件必须包含版权声明头部
2. 遵循 Go 代码规范和项目约定
3. 新增功能必须编写对应的单元测试
4. API 调用使用 `runtime.CallAPI` 方法
5. 输出格式化使用 `runtime.OutFormat` 和 `output.PrintTable`
