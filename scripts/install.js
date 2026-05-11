const os = require('os');
const path = require('path');
const fs = require('fs');
const https = require('https');
const { execSync } = require('child_process');

const platform = os.platform();
const arch = os.arch();
const version = require('../package.json').version;

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
