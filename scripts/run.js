#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');
const fs = require('fs');

// 获取二进制文件路径
const platform = process.platform;
const binaryName = platform === 'win32' ? 'soke-cli.exe' : 'soke-cli';
const binaryPath = path.join(__dirname, '..', 'bin', binaryName);

// 检查二进制文件是否存在
if (!fs.existsSync(binaryPath)) {
  console.error('错误: 未找到 soke-cli 二进制文件');
  console.error('请尝试重新安装: npm install -g @sokeai/cli');
  process.exit(1);
}

// 检查二进制文件是否可执行
try {
  fs.accessSync(binaryPath, fs.constants.X_OK);
} catch (err) {
  // 如果不可执行，尝试设置权限
  if (platform !== 'win32') {
    try {
      fs.chmodSync(binaryPath, 0o755);
    } catch (chmodErr) {
      console.error('错误: 无法设置二进制文件执行权限');
      process.exit(1);
    }
  }
}

// 执行二进制文件，传递所有命令行参数
const child = spawn(binaryPath, process.argv.slice(2), {
  stdio: 'inherit',
  env: process.env
});

// 处理子进程退出
child.on('exit', (code, signal) => {
  if (signal) {
    process.kill(process.pid, signal);
  } else {
    process.exit(code);
  }
});

// 处理错误
child.on('error', (err) => {
  console.error('执行 soke-cli 时出错:', err.message);
  process.exit(1);
});
