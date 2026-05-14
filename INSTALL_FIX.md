# install.js 修复说明

## 🐛 问题描述

用户执行 `npm install -g @sokeai/cli@latest` 后，soke-course skill 没有被自动安装和注册。

## 🔍 根本原因

在 `scripts/install.js` 文件中，`skillNames` 数组只包含了 `soke-shared` 和 `soke-exam`，**没有包含 `soke-course`**。

### 问题代码位置

**第 98 行**（syncSkillsToSokeclawWorkspace 函数）：
```javascript
const skillNames = ['soke-shared', 'soke-exam']; // ❌ 缺少 soke-course
```

**第 168 行**（syncSkillsToWorkclawRegistry 函数）：
```javascript
const skillNames = ['soke-shared', 'soke-exam']; // ❌ 缺少 soke-course
```

**第 175-209 行**（registry 注册）：
- ✅ 有 soke-exam 的注册信息
- ✅ 有 soke-shared 的注册信息
- ❌ **缺少 soke-course 的注册信息**

## ✅ 修复内容

### 1. 更新 skillNames 数组（2处）

**第 98 行**：
```javascript
const skillNames = ['soke-shared', 'soke-exam', 'soke-course']; // ✅ 添加 soke-course
```

**第 168 行**：
```javascript
const skillNames = ['soke-shared', 'soke-exam', 'soke-course']; // ✅ 添加 soke-course
```

### 2. 添加 soke-course 的注册信息

在第 210 行之前添加：
```javascript
upsertSkillRegistryEntry(registry, {
  id: 'skill:soke-course',
  name: 'soke-course',
  displayName: '授客课程管理',
  description: '授客课程管理：查询课程、课程分类、课程详情、学习记录。查询课程列表、课程分类、课程用户学习记录、课件列表、人脸识别记录。',
  source: { type: defaultSourceType, slug: '', url: '' },
  install: {
    path: path.join(workclawSkillInstallDir, 'soke-course'),
    installedAt: '',
    updatedAt: '',
    version: '1.0.0'
  },
  state: { enabled: true, health: 'ok', lastError: '' },
  runtime: { supported: ['openclaw'], enabled: ['openclaw'], primary: 'openclaw' },
  security: { riskLevel: 'normal', requiresApproval: false },
  metadata: { emoji: '📚', homepage: '', requires: { bins: ['soke-cli'] } }
});
```

## 📊 修复后的效果

当用户执行 `npm install -g @sokeai/cli@latest` 时，install.js 会自动：

1. ✅ 复制 `soke-course` 到 `~/.workclaw/skills/soke-course/`
2. ✅ 在 `~/.workclaw/skills/registry.json` 中注册 soke-course
3. ✅ 设置正确的元数据（emoji: 📚, requires: bins: ['soke-cli']）

## 🔄 安装流程

### install.js 的执行流程

```
npm install -g @sokeai/cli@latest
    ↓
package.json 的 postinstall 钩子
    ↓
scripts/install.js 执行
    ↓
1. 下载二进制文件
    ↓
2. syncSkillsToSokeclawWorkspace()
   - 复制 skills 到 ~/.sokeclaw/... 和 ~/.zev/...
   - skillNames: ['soke-shared', 'soke-exam', 'soke-course'] ✅
    ↓
3. syncSkillsToWorkclawRegistry()
   - 复制 skills 到 ~/.workclaw/skills/
   - skillNames: ['soke-shared', 'soke-exam', 'soke-course'] ✅
   - 注册到 registry.json
     * soke-exam ✅
     * soke-shared ✅
     * soke-course ✅ (新添加)
    ↓
4. 完成安装
```

## 📝 验证修改

### 验证 skillNames 数组
```bash
grep -n "skillNames = \[" scripts/install.js
```

**预期输出**：
```
98:  const skillNames = ['soke-shared', 'soke-exam', 'soke-course'];
168:  const skillNames = ['soke-shared', 'soke-exam', 'soke-course'];
```

### 验证注册信息
```bash
grep -n "skill:soke-course" scripts/install.js
```

**预期输出**：
```
212:    id: 'skill:soke-course',
```

## 🚀 下一步

### 1. 提交代码
```bash
git add scripts/install.js
git commit -m "fix: add soke-course to install.js skillNames and registry"
```

### 2. 更新版本号
```bash
# 编辑 package.json，将版本从 1.0.19 改为 1.0.20
npm version patch -m "fix: add soke-course skill to auto-install"
```

### 3. 发布新版本
```bash
# 推送代码
git push origin mycli:main
git push origin --tags

# 发布到 npm
npm publish --access public
```

### 4. 用户重新安装
用户执行以下命令即可获取修复：
```bash
npm cache clean --force
npm install -g @sokeai/cli@latest
```

## 📋 修改清单

| 文件 | 行号 | 修改内容 | 状态 |
|------|------|----------|------|
| scripts/install.js | 98 | 添加 'soke-course' 到 skillNames | ✅ 已修复 |
| scripts/install.js | 168 | 添加 'soke-course' 到 skillNames | ✅ 已修复 |
| scripts/install.js | 210+ | 添加 soke-course 注册信息 | ✅ 已修复 |

## 🎯 影响范围

### 修复前
- ❌ soke-course 不会被自动安装
- ❌ soke-course 不会被注册到 registry.json
- ❌ 用户需要手动复制和注册

### 修复后
- ✅ soke-course 会被自动安装到 ~/.workclaw/skills/
- ✅ soke-course 会被自动注册到 registry.json
- ✅ 用户无需任何手动操作

## 🔗 相关文件

- `/Users/edy/www/soke-lark-cli/soke-cli/scripts/install.js` - 安装脚本（已修复）
- `/Users/edy/www/soke-lark-cli/soke-cli/package.json` - 包配置
- `/Users/edy/www/soke-lark-cli/soke-cli/skills/soke-course/` - skill 文件
- `~/.workclaw/skills/registry.json` - skill 注册表

## 📚 相关文档

- [TROUBLESHOOTING.md](./TROUBLESHOOTING.md) - 问题排查文档
- [PUBLISH_GUIDE.md](./PUBLISH_GUIDE.md) - 发布指南

---

**修复日期**: 2024-05-14
**修复内容**: 添加 soke-course 到 install.js 的自动安装和注册流程
**影响版本**: 1.0.20+
**状态**: ✅ 已修复，待发布
