#!/usr/bin/env node

/**
 * 测试 install.js 的自动检测功能
 */

const fs = require('fs');
const path = require('path');

// 复制 detectSkillNames 函数
function detectSkillNames(packagedSkillsDir) {
  if (!fs.existsSync(packagedSkillsDir)) return [];

  try {
    const entries = fs.readdirSync(packagedSkillsDir, { withFileTypes: true });
    return entries
      .filter(entry => entry.isDirectory() && entry.name.startsWith('soke-'))
      .map(entry => entry.name)
      .sort();
  } catch (_) {
    return [];
  }
}

// 复制 parseSkillMetadata 函数
function parseSkillMetadata(skillDir) {
  const skillMdPath = path.join(skillDir, 'SKILL.md');
  if (!fs.existsSync(skillMdPath)) {
    return null;
  }

  try {
    const content = fs.readFileSync(skillMdPath, 'utf8');

    const frontmatterMatch = content.match(/^---\n([\s\S]*?)\n---/);
    if (!frontmatterMatch) return null;

    const frontmatter = frontmatterMatch[1];
    const metadata = {};

    const nameMatch = frontmatter.match(/^name:\s*(.+)$/m);
    if (nameMatch) metadata.name = nameMatch[1].trim();

    const summaryMatch = frontmatter.match(/^summary:\s*(.+)$/m);
    if (summaryMatch) metadata.summary = summaryMatch[1].trim();

    const descMatch = frontmatter.match(/^description:\s*["'](.+)["']$/m);
    if (descMatch) {
      metadata.description = descMatch[1].trim();
    } else {
      const descMatch2 = frontmatter.match(/^description:\s*(.+)$/m);
      if (descMatch2) metadata.description = descMatch2[1].trim();
    }

    const versionMatch = frontmatter.match(/^version:\s*(.+)$/m);
    if (versionMatch) metadata.version = versionMatch[1].trim();

    const binsMatch = frontmatter.match(/bins:\s*\[(.+?)\]/);
    if (binsMatch) {
      metadata.bins = binsMatch[1].split(',').map(b => b.trim().replace(/['"]/g, ''));
    }

    return metadata;
  } catch (_) {
    return null;
  }
}

// 测试
const packageRoot = path.join(__dirname, '..');
const skillsDir = path.join(packageRoot, 'skills');

console.log('🔍 测试自动检测功能\n');
console.log('Skills 目录:', skillsDir);
console.log('');

// 检测所有 skills
const skillNames = detectSkillNames(skillsDir);

console.log(`✅ 检测到 ${skillNames.length} 个 skills:\n`);

// 显示每个 skill 的详细信息
for (const skillName of skillNames) {
  const skillDir = path.join(skillsDir, skillName);
  const metadata = parseSkillMetadata(skillDir);

  console.log(`📦 ${skillName}`);

  if (metadata) {
    console.log(`   名称: ${metadata.name || '未知'}`);
    console.log(`   摘要: ${metadata.summary || '未知'}`);
    console.log(`   版本: ${metadata.version || '未知'}`);
    console.log(`   依赖: ${metadata.bins ? metadata.bins.join(', ') : '无'}`);
    console.log(`   描述: ${metadata.description ? metadata.description.substring(0, 60) + '...' : '未知'}`);
  } else {
    console.log('   ⚠️  无法解析 SKILL.md');
  }

  console.log('');
}

// 验证结果
console.log('📊 验证结果:\n');

const expectedSkills = ['soke-course', 'soke-exam', 'soke-shared'];
const missingSkills = expectedSkills.filter(s => !skillNames.includes(s));
const extraSkills = skillNames.filter(s => !expectedSkills.includes(s));

if (missingSkills.length > 0) {
  console.log(`❌ 缺少的 skills: ${missingSkills.join(', ')}`);
} else {
  console.log('✅ 所有预期的 skills 都已检测到');
}

if (extraSkills.length > 0) {
  console.log(`ℹ️  额外的 skills: ${extraSkills.join(', ')}`);
}

console.log('');
console.log('🎉 测试完成！');
console.log('');
console.log('💡 提示:');
console.log('   - 新增 skill 时，只需在 skills/ 目录下创建 soke-* 目录');
console.log('   - 确保每个 skill 都有 SKILL.md 文件，包含完整的 frontmatter');
console.log('   - install.js 会自动检测并注册所有 skills');
