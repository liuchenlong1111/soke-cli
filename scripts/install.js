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

function detectSokeclawWorkspaceSkillsDir() {
  const homeDir = os.homedir();
  const defaultSkillsDir = path.join(
    homeDir,
    '.sokeclaw',
    'openai-agents',
    'workspaces',
    'main',
    'skills'
  );

  const workclawConfigPath = path.join(homeDir, '.sokeclaw', 'workclaw.json');
  if (!fs.existsSync(workclawConfigPath)) return defaultSkillsDir;

  try {
    const configText = fs.readFileSync(workclawConfigPath, 'utf8');
    const config = JSON.parse(configText);
    const workspaceDir = config?.defaults?.agents?.openaiAgents?.main?.workspace;
    if (typeof workspaceDir === 'string' && workspaceDir.length > 0) {
      return path.join(workspaceDir, 'skills');
    }
  } catch (_) {}

  return defaultSkillsDir;
}

function syncSkillsToSokeclawWorkspace() {
  const packageRoot = path.join(__dirname, '..');
  const packagedSkillsDir = path.join(packageRoot, 'skills');
  if (!fs.existsSync(packagedSkillsDir)) return;

  const sokeclawSkillsDir = detectSokeclawWorkspaceSkillsDir();
  const sokeclawRootDir = path.join(os.homedir(), '.sokeclaw');
  if (!fs.existsSync(sokeclawRootDir)) return;

  try {
    fs.mkdirSync(sokeclawSkillsDir, { recursive: true });
  } catch (_) {
    return;
  }

  const skillNames = ['soke-shared', 'soke-exam'];
  for (const skillName of skillNames) {
    const src = path.join(packagedSkillsDir, skillName);
    const dest = path.join(sokeclawSkillsDir, skillName);
    if (fs.existsSync(src)) {
      copyDirRecursive(src, dest);
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

  const skillNames = ['soke-shared', 'soke-exam'];
  for (const skillName of skillNames) {
    const src = path.join(packagedSkillsDir, skillName);
    const dest = path.join(workclawSkillInstallDir, skillName);
    if (fs.existsSync(src)) copyDirRecursive(src, dest);
  }

  upsertSkillRegistryEntry(registry, {
    id: 'skill:soke-exam',
    name: 'soke-exam',
    displayName: '授客考试管理',
    description: '授客考试管理：查询考试、考试分类、考试用户成绩、考试详情。',
    source: { type: defaultSourceType, slug: '', url: '' },
    install: {
      path: path.join(workclawSkillInstallDir, 'soke-exam'),
      installedAt: '',
      updatedAt: '',
      version: '1.0.0'
    },
    state: { enabled: true, health: 'ok', lastError: '' },
    runtime: { supported: ['openclaw'], enabled: ['openclaw'], primary: 'openclaw' },
    security: { riskLevel: 'normal', requiresApproval: false },
    metadata: { emoji: '📝', homepage: '', requires: { bins: ['soke-cli'] } }
  });

  upsertSkillRegistryEntry(registry, {
    id: 'skill:soke-shared',
    name: 'soke-shared',
    displayName: 'soke-shared 共享规则',
    description: '授客CLI共享基础：配置、登录、权限管理、错误处理、安全规则。',
    source: { type: defaultSourceType, slug: '', url: '' },
    install: {
      path: path.join(workclawSkillInstallDir, 'soke-shared'),
      installedAt: '',
      updatedAt: '',
      version: '1.0.0'
    },
    state: { enabled: true, health: 'ok', lastError: '' },
    runtime: { supported: ['openclaw'], enabled: ['openclaw'], primary: 'openclaw' },
    security: { riskLevel: 'normal', requiresApproval: false },
    metadata: { emoji: '🔧', homepage: '', requires: {} }
  });

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
  if (!isLikelyInteractiveInstall()) return;

  const shimPath = detectSokeCliShimPath();
  if (!shimPath) return;

  const targets = getPreferredLinkTargetPaths();
  if (targets.length === 0) return;

  let selectedTarget = null;
  for (const t of targets) {
    const dir = path.dirname(t);
    if (fs.existsSync(dir)) {
      selectedTarget = t;
      break;
    }
  }
  if (!selectedTarget) selectedTarget = targets[0];

  const dir = path.dirname(selectedTarget);
  let dirWritable = false;
  try {
    fs.accessSync(dir, fs.constants.W_OK);
    dirWritable = true;
  } catch (_) {}

  const question = dirWritable
    ? `检测到你可能在 sokeclaw（GUI）里遇到 “command not found: soke-cli”。是否创建链接 ${selectedTarget} 指向 ${shimPath} 以便 GUI 可直接找到？(y/N) `
    : `检测到你可能在 sokeclaw（GUI）里遇到 “command not found: soke-cli”。是否输出一条需要 sudo 的命令来创建链接 ${selectedTarget} 指向 ${shimPath}？(y/N) `;

  const ok = await promptYesNo(question);
  if (!ok) return;

  if (dirWritable) {
    const res = ensureSymlink(selectedTarget, shimPath);
    if (res.ok) {
      console.log(`已配置：${selectedTarget} -> ${shimPath}`);
    } else if (res.reason === 'exists_non_symlink') {
      console.log(`跳过：${selectedTarget} 已存在且不是软链接。`);
    } else {
      console.log(`创建软链接失败：${res.reason}`);
    }
    return;
  }

  console.log(`请执行以下命令完成配置（需要 sudo）：`);
  console.log(`sudo ln -sf "${shimPath}" "${selectedTarget}"`);
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
