const os = require('os');
const path = require('path');
const fs = require('fs');
const https = require('https');
const { execSync } = require('child_process');
const readline = require('readline');

const platform = os.platform();
const arch = os.arch();
const version = require('../package.json').version;

function copyDirRecursive(srcDir, destDir) {
  if (!fs.existsSync(srcDir)) return;
  if (!fs.existsSync(destDir)) fs.mkdirSync(destDir, { recursive: true });

  const entries = fs.readdirSync(srcDir, { withFileTypes: true });
  for (const entry of entries) {
    const srcPath = path.join(srcDir, entry.name);
    const destPath = path.join(destDir, entry.name);

    if (entry.isDirectory()) {
      copyDirRecursive(srcPath, destPath);
      continue;
    }

    if (entry.isSymbolicLink()) {
      try {
        const linkTarget = fs.readlinkSync(srcPath);
        try {
          fs.unlinkSync(destPath);
        } catch (_) {}
        fs.symlinkSync(linkTarget, destPath);
      } catch (_) {}
      continue;
    }

    fs.copyFileSync(srcPath, destPath);
  }
}

function detectSokeclawWorkspaceSkillsDirs() {
  const homeDir = os.homedir();
  const dirs = [];

  // 1. Sokeclaw 默认工作区
  const defaultSokeclawDir = path.join(
    homeDir,
    '.sokeclaw',
    'openai-agents',
    'workspaces',
    'main',
    'skills'
  );

  // 2. Zev 默认工作区 (Sokeclaw 可能会生成这个)
  const defaultZevDir = path.join(
    homeDir,
    '.zev',
    'openai-agents',
    'workspaces',
    'main',
    'skills'
  );

  // 3. 从 workclaw.json 读取配置
  const workclawConfigPath = path.join(homeDir, '.sokeclaw', 'workclaw.json');
  let configuredDir = null;
  if (fs.existsSync(workclawConfigPath)) {
    try {
      const configText = fs.readFileSync(workclawConfigPath, 'utf8');
      const config = JSON.parse(configText);
      const workspaceDir = config?.defaults?.agents?.openaiAgents?.main?.workspace;
      if (typeof workspaceDir === 'string' && workspaceDir.length > 0) {
        configuredDir = path.join(workspaceDir, 'skills');
      }
    } catch (_) {}
  }

  // 收集所有有效的潜在目录
  if (configuredDir) {
    dirs.push(configuredDir);
  } else {
    dirs.push(defaultSokeclawDir);
  }
  
  // 始终将 .zev 目录也加入同步列表
  dirs.push(defaultZevDir);

  return dirs;
}

/**
 * 自动检测 skills 目录中的所有 skill
 * 只包含以 'soke-' 开头的目录
 */
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

function syncSkillsToSokeclawWorkspace() {
  const packageRoot = path.join(__dirname, '..');
  const packagedSkillsDir = path.join(packageRoot, 'skills');
  if (!fs.existsSync(packagedSkillsDir)) return;

  const targetDirs = detectSokeclawWorkspaceSkillsDirs();
  const skillNames = detectSkillNames(packagedSkillsDir);

  for (const targetDir of targetDirs) {
    try {
      fs.mkdirSync(targetDir, { recursive: true });
      for (const skillName of skillNames) {
        const src = path.join(packagedSkillsDir, skillName);
        const dest = path.join(targetDir, skillName);
        if (fs.existsSync(src)) {
          copyDirRecursive(src, dest);
        }
      }
    } catch (_) {
      // 忽略单个目录的写入失败
    }
  }
}

function upsertSkillRegistryEntry(registry, entry) {
  if (!registry || typeof registry !== 'object') return;
  if (!Array.isArray(registry.skills)) registry.skills = [];

  const idx = registry.skills.findIndex((s) => s && s.id === entry.id);
  if (idx >= 0) {
    registry.skills[idx] = { ...registry.skills[idx], ...entry };
    return;
  }

  registry.skills.push(entry);
}

/**
 * 从 SKILL.md 文件中解析元数据
 */
function parseSkillMetadata(skillDir) {
  const skillMdPath = path.join(skillDir, 'SKILL.md');
  if (!fs.existsSync(skillMdPath)) {
    return null;
  }

  try {
    const content = fs.readFileSync(skillMdPath, 'utf8');

    // 解析 frontmatter (YAML)
    const frontmatterMatch = content.match(/^---\n([\s\S]*?)\n---/);
    if (!frontmatterMatch) return null;

    const frontmatter = frontmatterMatch[1];
    const metadata = {};

    // 解析 name
    const nameMatch = frontmatter.match(/^name:\s*(.+)$/m);
    if (nameMatch) metadata.name = nameMatch[1].trim();

    // 解析 summary (用作 displayName)
    const summaryMatch = frontmatter.match(/^summary:\s*(.+)$/m);
    if (summaryMatch) metadata.summary = summaryMatch[1].trim();

    // 解析 description
    const descMatch = frontmatter.match(/^description:\s*["'](.+)["']$/m);
    if (descMatch) {
      metadata.description = descMatch[1].trim();
    } else {
      const descMatch2 = frontmatter.match(/^description:\s*(.+)$/m);
      if (descMatch2) metadata.description = descMatch2[1].trim();
    }

    // 解析 version
    const versionMatch = frontmatter.match(/^version:\s*(.+)$/m);
    if (versionMatch) metadata.version = versionMatch[1].trim();

    // 解析 metadata.requires.bins
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
 * 根据 skill 名称推断 emoji
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

function syncSkillsToWorkclawRegistry() {
  const homeDir = os.homedir();
  const workclawRootDir = path.join(homeDir, '.workclaw');
  const workclawSkillsDir = path.join(workclawRootDir, 'skills');
  const registryPath = path.join(workclawSkillsDir, 'registry.json');

  if (!fs.existsSync(workclawRootDir)) return;
  try {
    fs.mkdirSync(workclawSkillsDir, { recursive: true });
  } catch (_) {
    return;
  }

  const packageRoot = path.join(__dirname, '..');
  const packagedSkillsDir = path.join(packageRoot, 'skills');
  if (!fs.existsSync(packagedSkillsDir)) return;

  let registry;
  try {
    registry = JSON.parse(fs.readFileSync(registryPath, 'utf8'));
  } catch (_) {
    registry = { version: 1, migrations: {}, skills: [] };
  }

  if (registry.version == null) registry.version = 1;
  if (!registry.migrations) registry.migrations = {};

  const existingSkills = Array.isArray(registry.skills) ? registry.skills : [];
  const existingInstallPaths = existingSkills
    .map((s) => s?.install?.path)
    .filter((p) => typeof p === 'string');

  const workclawSkillInstallDir = workclawSkillsDir;

  const defaultSourceType =
    typeof existingSkills[0]?.source?.type === 'string' && existingSkills[0]?.source?.type
      ? existingSkills[0].source.type
      : 'local';

  // 自动检测所有 skills
  const skillNames = detectSkillNames(packagedSkillsDir);

  // 复制所有 skills
  for (const skillName of skillNames) {
    const src = path.join(packagedSkillsDir, skillName);
    const dest = path.join(workclawSkillInstallDir, skillName);
    if (fs.existsSync(src)) copyDirRecursive(src, dest);
  }

  // 自动注册所有 skills
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
      source: { type: defaultSourceType, slug: '', url: '' },
      install: {
        path: path.join(workclawSkillInstallDir, skillName),
        installedAt: '',
        updatedAt: '',
        version: version
      },
      state: { enabled: true, health: 'ok', lastError: '' },
      runtime: { supported: ['openclaw'], enabled: ['openclaw'], primary: 'openclaw' },
      security: { riskLevel: 'normal', requiresApproval: false },
      metadata: { emoji: emoji, homepage: '', requires: requires }
    });
  }

  try {
    fs.writeFileSync(registryPath, JSON.stringify(registry, null, 2));
  } catch (_) {}
}

function isLikelyInteractiveInstall() {
  if (!process.stdin.isTTY) return false;
  if (!process.stdout.isTTY) return false;
  if (process.env.CI) return false;
  if (process.env.npm_config_yes === 'true') return false;
  return true;
}

function promptYesNo(message) {
  return new Promise((resolve) => {
    const rl = readline.createInterface({
      input: process.stdin,
      output: process.stdout
    });

    rl.question(message, (answer) => {
      rl.close();
      const normalized = String(answer || '').trim().toLowerCase();
      resolve(normalized === 'y' || normalized === 'yes');
    });
  });
}

function detectNpmGlobalBinDir() {
  const prefix = process.env.npm_config_prefix;
  if (typeof prefix === 'string' && prefix.length > 0) {
    const binDir = platform === 'win32' ? prefix : path.join(prefix, 'bin');
    return binDir;
  }
  return null;
}

function detectSokeCliShimPath() {
  const binDir = detectNpmGlobalBinDir();
  if (!binDir) return null;

  const candidates =
    platform === 'win32'
      ? [path.join(binDir, 'soke-cli.cmd'), path.join(binDir, 'soke-cli.exe')]
      : [path.join(binDir, 'soke-cli')];

  for (const candidate of candidates) {
    if (fs.existsSync(candidate)) return candidate;
  }
  return null;
}

function getPreferredLinkTargetPaths() {
  if (platform === 'win32') return [];
  return ['/usr/local/bin/soke-cli', '/opt/homebrew/bin/soke-cli'];
}

function ensureSymlink(targetPath, sourcePath) {
  try {
    const existing = fs.lstatSync(targetPath);
    if (existing.isSymbolicLink()) {
      const currentTarget = fs.readlinkSync(targetPath);
      if (currentTarget === sourcePath) return { ok: true, changed: false };
      fs.unlinkSync(targetPath);
    } else {
      return { ok: false, changed: false, reason: 'exists_non_symlink' };
    }
  } catch (_) {}

  try {
    fs.symlinkSync(sourcePath, targetPath);
    return { ok: true, changed: true };
  } catch (err) {
    return { ok: false, changed: false, reason: err && err.message ? err.message : 'failed' };
  }
}

async function maybeAssistGuiPath() {
  const platform = os.platform();
  const binDir = path.join(__dirname, '..', 'bin');
  const binaryName = platform === 'win32' ? 'soke-cli.exe' : 'soke-cli';
  const binaryPath = path.join(binDir, binaryName);
  
  if (!fs.existsSync(binaryPath)) return Promise.resolve();

  // 我们不再弹窗询问，而是直接以静默模式尝试执行 setup-gui-env
  // 这会尝试创建软链，如果没权限，它会默默失败（不会报错打断安装）
  // 用户如果后面遇到 command not found，仍可以手动执行 soke-cli setup-gui-env
  return new Promise((resolve) => {
    try {
      const runJsPath = path.join(__dirname, 'run.js');
      const { spawn } = require('child_process');
      const child = spawn(process.execPath, [runJsPath, 'setup-gui-env', '--silent'], {
        stdio: 'ignore', // 静默执行，不污染 npm install 输出
        detached: true   // 允许子进程独立于父进程运行
      });
      // 断开与子进程的联系，让它在后台自己跑，不阻塞 npm 安装退出
      child.unref();
      resolve();
    } catch (e) {
      resolve();
    }
  });
}

// 平台映射
const platformMap = {
  'darwin': 'darwin',
  'linux': 'linux',
  'win32': 'windows'
};

const archMap = {
  'x64': 'amd64',
  'arm64': 'arm64'
};

const mappedPlatform = platformMap[platform];
const mappedArch = archMap[arch];

if (!mappedPlatform || !mappedArch) {
  console.error(`不支持的平台: ${platform}-${arch}`);
  console.error('支持的平台: darwin-x64, darwin-arm64, linux-x64, windows-x64');
  process.exit(1);
}

const binaryName = platform === 'win32' ? 'soke-cli.exe' : 'soke-cli';
const binaryFileName = `soke-cli-${mappedPlatform}-${mappedArch}${platform === 'win32' ? '.exe' : ''}`;

// GitHub Releases 下载地址
const downloadURL = `https://github.com/liuchenlong1111/soke-cli/releases/download/v${version}/${binaryFileName}`;

const binDir = path.join(__dirname, '..', 'bin');
const binaryPath = path.join(binDir, binaryName);

console.log(`正在为 ${mappedPlatform}-${mappedArch} 下载 soke-cli v${version}...`);
console.log(`下载地址: ${downloadURL}`);

// 创建 bin 目录
if (!fs.existsSync(binDir)) {
  fs.mkdirSync(binDir, { recursive: true });
}

// 下载文件
function downloadFile(url, dest) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(dest);

    https.get(url, (response) => {
      // 处理重定向
      if (response.statusCode === 301 || response.statusCode === 302) {
        file.close();
        fs.unlinkSync(dest);
        return downloadFile(response.headers.location, dest)
          .then(resolve)
          .catch(reject);
      }

      if (response.statusCode !== 200) {
        file.close();
        fs.unlinkSync(dest);
        return reject(new Error(`下载失败，HTTP 状态码: ${response.statusCode}`));
      }

      const totalSize = parseInt(response.headers['content-length'], 10);
      let downloadedSize = 0;
      let lastPercent = 0;

      response.on('data', (chunk) => {
        downloadedSize += chunk.length;
        const percent = Math.floor((downloadedSize / totalSize) * 100);
        if (percent !== lastPercent && percent % 10 === 0) {
          process.stdout.write(`\r下载进度: ${percent}%`);
          lastPercent = percent;
        }
      });

      response.pipe(file);

      file.on('finish', () => {
        file.close();
        console.log('\r下载完成!                    ');
        resolve();
      });
    }).on('error', (err) => {
      file.close();
      fs.unlinkSync(dest);
      reject(err);
    });
  });
}

// 执行下载
downloadFile(downloadURL, binaryPath)
  .then(() => {
    // 设置可执行权限（非 Windows 平台）
    if (platform !== 'win32') {
      try {
        fs.chmodSync(binaryPath, 0o755);
        console.log('已设置可执行权限');
      } catch (err) {
        console.error('设置可执行权限失败:', err.message);
        process.exit(1);
      }
    }

    try {
      syncSkillsToSokeclawWorkspace();
    } catch (_) {}

    try {
      syncSkillsToWorkclawRegistry();
    } catch (_) {}

    return maybeAssistGuiPath();
  })
  .then(() => {
    console.log('soke-cli 安装成功!');
    console.log(`二进制文件位置: ${binaryPath}`);
    console.log('\n使用方法:');
    console.log('  soke-cli --help');
    console.log('  soke-cli config init');
    console.log('  soke-cli auth login');
  })
  .catch((err) => {
    console.error('\n安装失败:', err.message);
    console.error('\n可能的原因:');
    console.error('1. 网络连接问题');
    console.error('2. GitHub Releases 中不存在该版本的二进制文件');
    console.error('3. 不支持当前平台');
    console.error('\n请访问 https://github.com/sokeai/soke-cli/releases 手动下载');
    process.exit(1);
  });
