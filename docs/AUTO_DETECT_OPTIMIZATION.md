# install.js 自动检测优化说明

## 🎯 优化目标

**问题**：每次新增 skill 时，需要手动修改 install.js 中的 `skillNames` 数组和注册信息，容易遗漏。

**解决方案**：实现自动检测和注册机制，新增 skill 时无需修改 install.js。

## ✨ 优化内容

### 1. 新增 `detectSkillNames()` 函数

**功能**：自动扫描 `skills/` 目录，检测所有以 `soke-` 开头的子目录。

```javascript
function detectSkillNames(packagedSkillsDir) {
  if (!fs.existsSync(packagedSkillsDir)) return [];

  try {
    const entries = fs.readdirSync(packagedSkillsDir, { withFileTypes: true });
    return entries
      .filter(entry => entry.isDirectory() && entry.name.startsWith('soke-'))
      .map(entry => entry.name)
      .sort(); // 排序确保一致性
  } catch (_) {
    return [];
  }
}
```

**特点**：
- ✅ 自动检测所有 `soke-*` 目录
- ✅ 按字母顺序排序
- ✅ 容错处理

### 2. 新增 `parseSkillMetadata()` 函数

**功能**：从 `SKILL.md` 文件的 frontmatter 中自动解析元数据。

```javascript
function parseSkillMetadata(skillDir) {
  const skillMdPath = path.join(skillDir, 'SKILL.md');
  if (!fs.existsSync(skillMdPath)) return null;

  try {
    const content = fs.readFileSync(skillMdPath, 'utf8');
    const frontmatterMatch = content.match(/^---\n([\s\S]*?)\n---/);
    if (!frontmatterMatch) return null;

    const frontmatter = frontmatterMatch[1];
    const metadata = {};

    // 解析 name, summary, description, version, bins
    // ...

    return metadata;
  } catch (_) {
    return null;
  }
}
```

**解析的字段**：
- `name` - skill 名称
- `summary` - 摘要（用作 displayName）
- `description` - 详细描述
- `version` - 版本号
- `metadata.requires.bins` - 依赖的二进制文件

### 3. 新增 `inferSkillEmoji()` 函数

**功能**：根据 skill 名称推断合适的 emoji。

```javascript
function inferSkillEmoji(skillName) {
  const emojiMap = {
    'soke-exam': '📝',
    'soke-course': '📚',
    'soke-shared': '🔧',
    'soke-user': '👤',
    'soke-contact': '📇',
    'soke-department': '🏢',
    'soke-approval': '✅',
    'soke-attendance': '📅',
    'soke-report': '📊'
  };
  return emojiMap[skillName] || '📦';
}
```

**特点**：
- ✅ 预定义常见 skill 的 emoji
- ✅ 未知 skill 使用默认 emoji 📦
- ✅ 易于扩展

### 4. 优化 `syncSkillsToSokeclawWorkspace()`

**修改前**：
```javascript
const skillNames = ['soke-shared', 'soke-exam', 'soke-course']; // ❌ 硬编码
```

**修改后**：
```javascript
const skillNames = detectSkillNames(packagedSkillsDir); // ✅ 自动检测
```

### 5. 优化 `syncSkillsToWorkclawRegistry()`

**修改前**：
```javascript
// ❌ 硬编码 skillNames
const skillNames = ['soke-shared', 'soke-exam', 'soke-course'];

// ❌ 手动注册每个 skill
upsertSkillRegistryEntry(registry, {
  id: 'skill:soke-exam',
  name: 'soke-exam',
  displayName: '授客考试管理',
  // ... 大量硬编码的配置
});
// ... 重复代码
```

**修改后**：
```javascript
// ✅ 自动检测所有 skills
const skillNames = detectSkillNames(packagedSkillsDir);

// ✅ 自动注册所有 skills
for (const skillName of skillNames) {
  const skillDir = path.join(packagedSkillsDir, skillName);
  const metadata = parseSkillMetadata(skillDir);

  if (!metadata || !metadata.name) {
    console.warn(`警告: 无法解析 ${skillName} 的元数据，跳过注册`);
    continue;
  }

  const displayName = metadata.summary || metadata.name;
  const description = metadata.description || `${displayName} - 授客AI CLI工具`;
  const version = metadata.version || '1.0.0';
  const emoji = inferSkillEmoji(skillName);
  const requires = metadata.bins ? { bins: metadata.bins } : {};

  upsertSkillRegistryEntry(registry, {
    id: `skill:${skillName}`,
    name: skillName,
    displayName: displayName,
    description: description,
    // ... 自动生成的配置
  });
}
```

## 📊 优化效果对比

### 修改前（手动维护）

| 操作 | 步骤 | 容易出错 |
|------|------|----------|
| 新增 skill | 1. 创建 skill 目录<br>2. 修改 install.js 第 98 行<br>3. 修改 install.js 第 168 行<br>4. 添加注册信息（~20行代码） | ❌ 容易遗漏 |
| 修改 skill 信息 | 1. 修改 SKILL.md<br>2. 修改 install.js 注册信息 | ❌ 需要两处修改 |
| 删除 skill | 1. 删除 skill 目录<br>2. 修改 install.js（3处） | ❌ 容易遗漏 |

### 修改后（自动检测）

| 操作 | 步骤 | 容易出错 |
|------|------|----------|
| 新增 skill | 1. 创建 skill 目录<br>2. 创建 SKILL.md（包含 frontmatter） | ✅ 自动检测 |
| 修改 skill 信息 | 1. 修改 SKILL.md | ✅ 自动同步 |
| 删除 skill | 1. 删除 skill 目录 | ✅ 自动移除 |

## 🎓 使用指南

### 新增 Skill 的标准流程

#### 1. 创建 skill 目录

```bash
mkdir -p skills/soke-newskill
```

#### 2. 创建 SKILL.md 文件

```markdown
---
name: soke-newskill
summary: 新功能管理（功能描述），通过 soke-cli 查询
version: 1.0.0
description: "新功能管理：详细描述功能。"
metadata:
  requires:
    bins: ["soke-cli"]
  cliHelp: "soke-cli newskill --help"
---

# 新功能管理 (newskill)

...
```

#### 3. 测试自动检测

```bash
node scripts/test-auto-detect.js
```

**预期输出**：
```
✅ 检测到 4 个 skills:

📦 soke-course
   ...

📦 soke-exam
   ...

📦 soke-newskill
   名称: soke-newskill
   摘要: 新功能管理（功能描述），通过 soke-cli 查询
   版本: 1.0.0
   依赖: soke-cli
   ...
```

#### 4. 发布

```bash
npm version patch -m "feat: add soke-newskill"
npm publish --access public
```

**就这么简单！无需修改 install.js！**

## ✅ 验证测试

### 测试脚本

创建了 `scripts/test-auto-detect.js` 用于测试自动检测功能。

```bash
node scripts/test-auto-detect.js
```

### 测试结果

```
🔍 测试自动检测功能

Skills 目录: /Users/edy/www/soke-lark-cli/soke-cli/skills

✅ 检测到 3 个 skills:

📦 soke-course
   名称: soke-course
   摘要: 授客课程管理（课程列表/分类/课程详情/学习记录），通过 soke-cli 查询
   版本: 1.0.0
   依赖: soke-cli

📦 soke-exam
   名称: soke-exam
   摘要: 授客考试管理（考试列表/分类/考试用户成绩/详情），通过 soke-cli 查询
   版本: 1.0.0
   依赖: soke-cli

📦 soke-shared
   名称: soke-shared
   摘要: 授客CLI共享基础（配置/登录/权限/错误处理/安全规则）
   版本: 1.0.0
   依赖: 无

📊 验证结果:

✅ 所有预期的 skills 都已检测到

🎉 测试完成！
```

## 📝 SKILL.md Frontmatter 规范

为了确保自动检测正常工作，每个 skill 的 SKILL.md 必须包含以下 frontmatter：

```yaml
---
name: soke-skillname              # 必需：skill 名称
summary: 简短描述                  # 必需：用作 displayName
version: 1.0.0                    # 必需：版本号
description: "详细描述"            # 必需：详细说明
metadata:
  requires:
    bins: ["soke-cli"]            # 可选：依赖的二进制文件
  cliHelp: "soke-cli cmd --help"  # 可选：帮助命令
---
```

### 字段说明

| 字段 | 必需 | 说明 | 示例 |
|------|------|------|------|
| name | ✅ | skill 名称，必须与目录名一致 | `soke-course` |
| summary | ✅ | 简短描述，用作 displayName | `授客课程管理（课程列表/分类/课程详情/学习记录），通过 soke-cli 查询` |
| version | ✅ | 版本号 | `1.0.0` |
| description | ✅ | 详细描述 | `授客课程管理：查询课程、课程分类、课程详情、学习记录。` |
| metadata.requires.bins | ❌ | 依赖的二进制文件 | `["soke-cli"]` |
| metadata.cliHelp | ❌ | 帮助命令 | `soke-cli course --help` |

## 🚀 优势总结

### 1. 零维护成本
- ✅ 新增 skill 无需修改 install.js
- ✅ 删除 skill 无需修改 install.js
- ✅ 修改 skill 信息只需修改 SKILL.md

### 2. 避免人为错误
- ✅ 不会忘记添加到 skillNames 数组
- ✅ 不会忘记添加注册信息
- ✅ 不会出现拼写错误

### 3. 保持一致性
- ✅ 所有 skill 的注册信息格式统一
- ✅ 元数据来源统一（SKILL.md）
- ✅ 易于维护和扩展

### 4. 提高开发效率
- ✅ 新增 skill 只需 2 步（创建目录 + 创建 SKILL.md）
- ✅ 测试脚本快速验证
- ✅ 自动化程度高

## 📋 修改清单

| 文件 | 修改内容 | 行数 |
|------|----------|------|
| scripts/install.js | 新增 `detectSkillNames()` 函数 | +15 |
| scripts/install.js | 新增 `parseSkillMetadata()` 函数 | +50 |
| scripts/install.js | 新增 `inferSkillEmoji()` 函数 | +15 |
| scripts/install.js | 优化 `syncSkillsToSokeclawWorkspace()` | -3, +3 |
| scripts/install.js | 优化 `syncSkillsToWorkclawRegistry()` | -80, +40 |
| scripts/test-auto-detect.js | 新增测试脚本 | +120 |

**总计**：
- 删除代码：~83 行（硬编码的 skillNames 和注册信息）
- 新增代码：~123 行（自动检测和解析逻辑）
- 净增加：~40 行
- 新增文件：1 个（test-auto-detect.js）

## 🎯 未来扩展

### 1. 支持更多元数据字段

可以在 `parseSkillMetadata()` 中添加更多字段解析：
- `author` - 作者
- `homepage` - 主页
- `tags` - 标签
- `category` - 分类

### 2. 支持自定义 emoji

在 SKILL.md frontmatter 中添加 `emoji` 字段：
```yaml
metadata:
  emoji: "🎓"
```

### 3. 验证 SKILL.md 格式

添加验证脚本，确保所有 SKILL.md 格式正确：
```bash
node scripts/validate-skills.js
```

### 4. 生成 Skills 文档

自动生成 skills 列表文档：
```bash
node scripts/generate-skills-doc.js > SKILLS.md
```

## 📚 相关文档

- [INSTALL_FIX.md](./INSTALL_FIX.md) - 之前的手动修复说明
- [TROUBLESHOOTING.md](./TROUBLESHOOTING.md) - 问题排查指南
- [PUBLISH_GUIDE.md](./PUBLISH_GUIDE.md) - 发布流程指南

---

**优化日期**: 2024-05-14
**优化内容**: 实现 skill 自动检测和注册机制
**影响版本**: 1.0.20+
**状态**: ✅ 已完成，待测试和发布
**维护成本**: 从"每次新增 skill 需修改 3 处"降低到"零维护"
