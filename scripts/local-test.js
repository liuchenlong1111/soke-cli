#!/usr/bin/env node

/**
 * 本地测试脚本 - 将 skills 分发到本地 AI Agent 环境
 * 用于开发完 skill 后在本地测试，无需发布到 npm
 *
 * 使用方法:
 *   node scripts/local-test.js
 *   node scripts/local-test.js --skill soke-course
 *   node scripts/local-test.js --clean
 */

const fs = require('fs');
const path = require('path');
const os = require('os');

// 颜色输出
const colors = {
  reset: '\x1b[0m',
  bright: '\x1b[1m',
  red: '\x1b[31m',
  green: '\x1b[32m',
  yellow: '\x1b[33m',
  blue: '\x1b[34m',
  cyan: '\x1b[36m'
};

function log(message, color = 'reset') {
  console.log(`${colors[color]}${message}${colors.reset}`);
}

function logSection(title) {
  console.log('');
  log(`${'='.repeat(60)}`, 'blue');
  log(`  ${title}`, 'bright');
  log(`${'='.repeat(60)}`, 'blue');
  console.log('');
}

function logSuccess(message) {
  log(`✅ ${message}`, 'green');
}

function logError(message) {
  log(`❌ ${message}`, 'red');
}

function logWarning(message) {
  log(`⚠️  ${message}`, 'yellow');
}

function logInfo(message) {
  log(`ℹ️  ${message}`, 'cyan');
}

/**
 * 递归复制目录
 */
function copyDirRecursive(srcDir, destDir) {
  if (!fs.existsSync(srcDir)) return false;

  if (!fs.existsSync(destDir)) {
    fs.mkdirSync(destDir, { recursive: true });
  }

  const entries = fs.readdirSync(srcDir, { withFileTypes: true });

  for (const entry of entries) {
    const srcPath = path.join(srcDir, entry.name);
    const destPath = path.join(destDir, entry.name);

    if (entry.isDirectory()) {
      copyDirRecursive(srcPath, destPath);
    } else if (entry.isSymbolicLink()) {
      try {
        const linkTarget = fs.readlinkSync(srcPath);
        try {
          fs.unlinkSync(destPath);
        } catch (_) {}
        fs.symlinkSync(linkTarget, destPath);
      } catch (_) {}
    } else {
      fs.copyFileSync(srcPath, destPath);
    }
  }

  return true;
}

/**
 * 检测所有 soke-* skills
 */
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

/**
 * 从 SKILL.md 解析元数据
 */
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

/**
 * 推断 skill emoji
 */
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

/**
 * 检测本地 AI Agent 环境
 * 只检测 agent 相关的目录，不包括全局安装目录
 */
function detectLocalAgentDirs() {
  const homeDir = os.homedir();
  const dirs = [];

  // 1. workclaw (Claude Code) - Agent skills 目录
  const workclawDir = path.join(homeDir, '.workclaw', 'skills');
  if (fs.existsSync(path.join(homeDir, '.workclaw'))) {
    dirs.push({
      name: 'workclaw (Claude Code)',
      path: workclawDir,
      registryPath: path.join(workclawDir, 'registry.json'),
      type: 'workclaw'
    });
  }

  // 2. claude (Claude Desktop) - Agent skills 目录
  const claudeDir = path.join(homeDir, '.claude', 'skills');
  if (fs.existsSync(path.join(homeDir, '.claude'))) {
    dirs.push({
      name: 'claude (Claude Desktop)',
      path: claudeDir,
      registryPath: null,
      type: 'claude'
    });
  }

  // 3. sokeclaw - Agent skills 目录
  const sokeclawDir = path.join(homeDir, '.sokeclaw', 'openai-agents', 'workspaces', 'main', 'skills');
  if (fs.existsSync(path.join(homeDir, '.sokeclaw'))) {
    dirs.push({
      name: 'sokeclaw',
      path: sokeclawDir,
      registryPath: null,
      type: 'sokeclaw'
    });
  }

  // 4. zev - Agent skills 目录
  const zevDir = path.join(homeDir, '.zev', 'openai-agents', 'workspaces', 'main', 'skills');
  if (fs.existsSync(path.join(homeDir, '.zev'))) {
    dirs.push({
      name: 'zev',
      path: zevDir,
      registryPath: null,
      type: 'zev'
    });
  }

  return dirs;
}

/**
 * 更新 workclaw registry.json
 */
function updateWorkclawRegistry(registryPath, skillName, metadata, skillInstallPath) {
  let registry;

  try {
    registry = JSON.parse(fs.readFileSync(registryPath, 'utf8'));
  } catch (_) {
    registry = { version: 1, migrations: {}, skills: [] };
  }

  if (registry.version == null) registry.version = 1;
  if (!registry.migrations) registry.migrations = {};
  if (!Array.isArray(registry.skills)) registry.skills = [];

  const displayName = metadata.summary || metadata.name || skillName;
  const description = metadata.description || `${displayName} - 授客AI CLI工具`;
  const version = metadata.version || '1.0.0';
  const emoji = inferSkillEmoji(skillName);
  const requires = metadata.bins ? { bins: metadata.bins } : {};

  const skillEntry = {
    id: `skill:${skillName}`,
    name: skillName,
    displayName: displayName,
    description: description,
    source: { type: 'local', slug: '', url: '' },
    install: {
      path: skillInstallPath,
      installedAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      version: version
    },
    state: { enabled: true, health: 'ok', lastError: '' },
    runtime: { supported: ['openclaw'], enabled: ['openclaw'], primary: 'openclaw' },
    security: { riskLevel: 'normal', requiresApproval: false },
    metadata: { emoji: emoji, homepage: '', requires: requires }
  };

  const idx = registry.skills.findIndex(s => s && s.id === skillEntry.id);
  if (idx >= 0) {
    registry.skills[idx] = { ...registry.skills[idx], ...skillEntry };
  } else {
    registry.skills.push(skillEntry);
  }

  fs.writeFileSync(registryPath, JSON.stringify(registry, null, 2));
}

/**
 * 分发单个 skill 到所有本地环境
 */
function distributeSkill(skillName, packagedSkillsDir, targetDirs) {
  const srcDir = path.join(packagedSkillsDir, skillName);

  if (!fs.existsSync(srcDir)) {
    logError(`Skill 目录不存在: ${srcDir}`);
    return false;
  }

  const metadata = parseSkillMetadata(srcDir);
  if (!metadata || !metadata.name) {
    logWarning(`无法解析 ${skillName} 的元数据，将使用默认值`);
  }

  logInfo(`分发 ${skillName}...`);
  console.log('');

  let successCount = 0;

  for (const target of targetDirs) {
    try {
      fs.mkdirSync(target.path, { recursive: true });

      const destDir = path.join(target.path, skillName);
      const success = copyDirRecursive(srcDir, destDir);

      if (success) {
        logSuccess(`  → ${target.name}`);
        logInfo(`     ${destDir}`);

        // 更新 workclaw registry
        if (target.type === 'workclaw' && target.registryPath) {
          updateWorkclawRegistry(target.registryPath, skillName, metadata, destDir);
          logInfo(`     已更新 registry.json`);
        }

        successCount++;
      } else {
        logError(`  → ${target.name} (复制失败)`);
      }
    } catch (err) {
      logError(`  → ${target.name} (错误: ${err.message})`);
    }
    console.log('');
  }

  return successCount > 0;
}

/**
 * 清理所有本地环境中的 skills
 */
function cleanSkills(targetDirs, skillNames) {
  logSection('清理本地 Skills');

  for (const target of targetDirs) {
    logInfo(`清理 ${target.name}...`);

    for (const skillName of skillNames) {
      const skillDir = path.join(target.path, skillName);

      if (fs.existsSync(skillDir)) {
        try {
          fs.rmSync(skillDir, { recursive: true, force: true });
          logSuccess(`  ✓ 删除 ${skillName}`);
        } catch (err) {
          logError(`  ✗ 删除 ${skillName} 失败: ${err.message}`);
        }
      }
    }

    // 清理 workclaw registry
    if (target.type === 'workclaw' && target.registryPath && fs.existsSync(target.registryPath)) {
      try {
        const registry = JSON.parse(fs.readFileSync(target.registryPath, 'utf8'));
        if (Array.isArray(registry.skills)) {
          const before = registry.skills.length;
          registry.skills = registry.skills.filter(s => {
            return !s || !s.name || !skillNames.includes(s.name);
          });
          const after = registry.skills.length;

          if (before !== after) {
            fs.writeFileSync(target.registryPath, JSON.stringify(registry, null, 2));
            logSuccess(`  ✓ 清理 registry.json (删除 ${before - after} 个条目)`);
          }
        }
      } catch (err) {
        logError(`  ✗ 清理 registry.json 失败: ${err.message}`);
      }
    }

    console.log('');
  }
}

/**
 * 主函数
 */
function main() {
  const args = process.argv.slice(2);
  const isClean = args.includes('--clean');
  const skillArg = args.find(arg => arg.startsWith('--skill='));
  const specificSkill = skillArg ? skillArg.split('=')[1] : null;

  logSection('本地 Skill 测试工具');

  // 检测项目目录
  const packageRoot = path.join(__dirname, '..');
  const packagedSkillsDir = path.join(packageRoot, 'skills');

  if (!fs.existsSync(packagedSkillsDir)) {
    logError(`Skills 目录不存在: ${packagedSkillsDir}`);
    process.exit(1);
  }

  // 检测所有 skills
  const allSkills = detectSkillNames(packagedSkillsDir);

  if (allSkills.length === 0) {
    logError('未检测到任何 soke-* skills');
    process.exit(1);
  }

  logInfo(`检测到 ${allSkills.length} 个 skills: ${allSkills.join(', ')}`);
  console.log('');

  // 检测本地 AI Agent 环境
  const targetDirs = detectLocalAgentDirs();

  if (targetDirs.length === 0) {
    logError('未检测到任何本地 AI Agent 环境');
    logInfo('支持的环境:');
    logInfo('  • workclaw (~/.workclaw/skills/)');
    logInfo('  • claude (~/.claude/skills/)');
    logInfo('  • sokeclaw (~/.sokeclaw/openai-agents/workspaces/main/skills/)');
    logInfo('  • zev (~/.zev/openai-agents/workspaces/main/skills/)');
    console.log('');
    logWarning('注意: 此脚本只分发到 Agent skills 目录，不包括全局安装目录');
    process.exit(1);
  }

  logInfo(`检测到 ${targetDirs.length} 个本地环境:`);
  for (const target of targetDirs) {
    logInfo(`  • ${target.name}`);
    logInfo(`    ${target.path}`);
  }
  console.log('');

  // 清理模式
  if (isClean) {
    cleanSkills(targetDirs, allSkills);
    logSuccess('清理完成！');
    return;
  }

  // 确定要分发的 skills
  const skillsToDistribute = specificSkill
    ? (allSkills.includes(specificSkill) ? [specificSkill] : [])
    : allSkills;

  if (skillsToDistribute.length === 0) {
    logError(`Skill 不存在: ${specificSkill}`);
    process.exit(1);
  }

  // 分发 skills
  logSection('分发 Skills 到本地环境');

  let totalSuccess = 0;
  let totalFailed = 0;

  for (const skillName of skillsToDistribute) {
    const success = distributeSkill(skillName, packagedSkillsDir, targetDirs);
    if (success) {
      totalSuccess++;
    } else {
      totalFailed++;
    }
  }

  // 总结
  logSection('分发完成');

  logInfo(`总计: ${skillsToDistribute.length} 个 skills`);
  logSuccess(`成功: ${totalSuccess} 个`);
  if (totalFailed > 0) {
    logError(`失败: ${totalFailed} 个`);
  }
  console.log('');

  // 下一步提示
  logSection('下一步');
  console.log('');
  log('1. 重启你的 AI Agent', 'cyan');
  log('   • Claude Code: 重启 VS Code 或重新打开 Claude Code 窗口', 'cyan');
  log('   • Claude Desktop: 重启 Claude Desktop 应用', 'cyan');
  log('   • Sokeclaw: 重启 sokeclaw 进程', 'cyan');
  log('   • Zev: 重启 zev 进程', 'cyan');
  console.log('');
  log('2. 在对话中测试 skill 功能', 'cyan');
  log('   例如: "查询课程列表" 或 "查询考试成绩"', 'cyan');
  console.log('');
  log('3. 验证命令是否可用:', 'cyan');
  console.log('');
  for (const skillName of skillsToDistribute) {
    log(`   soke-cli ${skillName.replace('soke-', '')} --help`, 'yellow');
  }
  console.log('');
  log('4. 测试完成后，可以清理:', 'cyan');
  log('   node scripts/local-test.js --clean', 'yellow');
  console.log('');
  log('💡 提示:', 'cyan');
  log('   • 此脚本只分发到 Agent skills 目录', 'cyan');
  log('   • 全局安装 (npm install -g) 需要单独处理', 'cyan');
  log('   • 修改后重新运行此脚本即可更新', 'cyan');
  console.log('');
}

// 运行
try {
  main();
} catch (err) {
  logError(`发生错误: ${err.message}`);
  console.error(err);
  process.exit(1);
}
