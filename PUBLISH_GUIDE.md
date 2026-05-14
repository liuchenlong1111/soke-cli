# 发布 soke-course Skill 到 NPM

## 📋 发布前检查清单

- ✅ soke-course skill 文件已创建（5个文件）
- ✅ 文件位于 `skills/soke-course/` 目录
- ✅ package.json 中包含 `skills/` 目录
- ⏳ 需要更新版本号
- ⏳ 需要提交代码
- ⏳ 需要发布到 npm

## 🚀 发布步骤

### 步骤1: 更新版本号

当前版本: `1.0.18`
建议新版本: `1.0.19`（添加了新的 soke-course skill）

```bash
cd /Users/edy/www/soke-lark-cli/soke-cli

# 方法1: 手动编辑 package.json
# 将 "version": "1.0.18" 改为 "version": "1.0.19"

# 方法2: 使用 npm version 命令
npm version patch -m "feat: add soke-course skill"
```

### 步骤2: 提交代码到 Git

```bash
# 如果还没有提交 soke-course 文件
git add skills/soke-course/
git commit -m "feat: add soke-course skill with complete documentation"

# 推送到远程仓库
git push origin mycli
```

### 步骤3: 发布到 NPM

#### 方法A: 使用自动化脚本（推荐）

```bash
cd /Users/edy/www/soke-lark-cli/soke-cli

# 运行发布脚本（会自动完成所有步骤）
./scripts/release.sh
```

这个脚本会自动：
1. ✅ 检查工具和 Git 状态
2. ✅ 运行测试
3. ✅ 编译所有平台二进制文件
4. ✅ 创建 Git 标签并推送
5. ✅ 上传到 GitHub Releases
6. ✅ 发布到 NPM

#### 方法B: 手动发布

```bash
cd /Users/edy/www/soke-lark-cli/soke-cli

# 1. 确保已登录 npm
npm whoami
# 如果未登录，运行: npm login

# 2. 编译二进制文件（如果需要）
./scripts/build-binaries.sh

# 3. 创建 Git 标签
git tag -a "v1.0.19" -m "Release v1.0.19 - Add soke-course skill"
git push origin v1.0.19

# 4. 发布到 npm
npm publish --access public

# 5. 创建 GitHub Release（可选）
gh release create "v1.0.19" \
  --title "v1.0.19 - Add soke-course skill" \
  --notes "新增 soke-course skill，支持课程查询、学习记录等功能"
```

## 📦 发布后验证

### 验证 NPM 包

```bash
# 查看最新版本
npm view @sokeai/cli version

# 安装最新版本
npm install -g @sokeai/cli@latest

# 验证 skill 是否存在
ls -la $(npm root -g)/@sokeai/cli/skills/soke-course/
```

### 验证 Skill 可用性

```bash
# 方法1: 直接使用命令
soke-cli course --help

# 方法2: 通过 npx skills 安装
npx skills add liuchenlong1111/soke-cli -y -g

# 验证 skill 文件
ls -la ~/.claude/skills/soke-course/
```

## ⚠️ 注意事项

### 1. 版本号规范

遵循语义化版本（Semantic Versioning）：
- **MAJOR.MINOR.PATCH** (例如: 1.0.19)
- **PATCH**: 修复 bug 或小改动（1.0.18 → 1.0.19）
- **MINOR**: 添加新功能（1.0.19 → 1.1.0）
- **MAJOR**: 破坏性更改（1.1.0 → 2.0.0）

建议：添加新 skill 使用 PATCH 或 MINOR 版本。

### 2. 发布前测试

```bash
# 测试命令是否可用
./soke-cli course --help
./soke-cli course +list-courses --help

# 测试 skill 文档
cat skills/soke-course/SKILL.md
```

### 3. NPM 发布权限

确保你有 `@sokeai/cli` 包的发布权限：
```bash
npm whoami
# 应该显示你的 npm 用户名

npm owner ls @sokeai/cli
# 应该包含你的用户名
```

### 4. GitHub 权限

如果使用 `gh` 命令创建 Release，需要：
```bash
gh auth status
# 确保已登录 GitHub
```

## 🔍 常见问题

### Q1: npm publish 失败，提示权限不足

**解决方案**:
```bash
# 重新登录
npm logout
npm login

# 或者联系包的所有者添加你为协作者
```

### Q2: Git 标签已存在

**解决方案**:
```bash
# 删除本地标签
git tag -d v1.0.19

# 删除远程标签
git push origin :refs/tags/v1.0.19

# 重新创建标签
git tag -a "v1.0.19" -m "Release v1.0.19"
git push origin v1.0.19
```

### Q3: 用户安装后看不到新 skill

**原因**: 用户可能安装了旧版本的缓存

**解决方案**:
```bash
# 用户端执行
npm cache clean --force
npm install -g @sokeai/cli@latest

# 或者指定版本
npm install -g @sokeai/cli@1.0.19
```

### Q4: npx skills add 没有获取到新 skill

**原因**: GitHub 仓库可能没有推送最新代码

**解决方案**:
```bash
# 确保推送到 GitHub
git push origin mycli

# 或者推送到 main 分支
git push origin mycli:main
```

## 📝 发布检查清单

发布前请确认：

- [ ] 版本号已更新（package.json）
- [ ] 代码已提交到 Git
- [ ] 代码已推送到 GitHub
- [ ] 已登录 npm（npm whoami）
- [ ] 测试命令可用（soke-cli course --help）
- [ ] skill 文档完整（SKILL.md, README.md 等）
- [ ] package.json 的 files 字段包含 skills/

发布后请验证：

- [ ] npm 上的版本已更新（npm view @sokeai/cli version）
- [ ] GitHub 上有对应的 tag 和 Release
- [ ] 用户可以安装最新版本（npm install -g @sokeai/cli@latest）
- [ ] skill 文件存在于安装目录
- [ ] 命令可以正常使用

## 🎯 快速发布命令

如果一切准备就绪，可以使用以下命令快速发布：

```bash
cd /Users/edy/www/soke-lark-cli/soke-cli

# 1. 更新版本号并提交
npm version patch -m "feat: add soke-course skill"

# 2. 推送代码和标签
git push origin mycli
git push origin --tags

# 3. 发布到 npm
npm publish --access public

# 4. 验证
npm view @sokeai/cli version
```

## 📚 相关文档

- [NPM 发布文档](https://docs.npmjs.com/cli/v8/commands/npm-publish)
- [语义化版本](https://semver.org/lang/zh-CN/)
- [GitHub Releases](https://docs.github.com/en/repositories/releasing-projects-on-github)

---

**准备发布**: v1.0.19
**新增内容**: soke-course skill（课程查询功能）
**发布日期**: 2024-05-14
