#!/usr/bin/env node

const { spawn, execSync } = require('child_process');
const path = require('path');
const fs = require('fs');

// 特殊命令拦截：修复 sokeclaw 的 GUI PATH 问题
if (process.argv[2] === 'setup-gui-env') {
  console.log('正在为 SokeClaw GUI 注入环境变量 PATH...');
  try {
    // 获取 node 和当前脚本的绝对路径目录，这通常也是 soke-cli 所在的 bin 目录
    const nodeDir = path.dirname(process.execPath);
    const cliDir = path.dirname(process.argv[1]);
    
    let targetPath = '';
    if (process.platform === 'darwin') {
      // 提取现有的系统 PATH
      const sysPath = execSync('sysctl -n getenv PATH 2>/dev/null || echo ""', { encoding: 'utf8' }).trim();
      const currentPath = process.env.PATH || '';
      
      // 合并并去重 PATH
      const pathSet = new Set([nodeDir, cliDir, ...currentPath.split(':'), ...sysPath.split(':')].filter(Boolean));
      targetPath = Array.from(pathSet).join(':');
      
      console.log(`注入的 PATH: ${targetPath.substring(0, 100)}...`);
      // launchctl setenv 在新版 macOS 中如果不在 login context 下可能提示 Not privileged，
      // 我们改用软链到 GUI 默认 PATH /usr/local/bin 的方式来解决（如果不可写则提示 sudo）
      const binSource = path.join(nodeDir, 'soke-cli');
      const binTarget = '/usr/local/bin/soke-cli';
      
      if (fs.existsSync(binSource)) {
        try {
          // 尝试创建软链
          if (fs.existsSync(binTarget)) {
            try { fs.unlinkSync(binTarget); } catch(e) {}
          }
          fs.symlinkSync(binSource, binTarget);
          console.log(`✅ 成功创建软链: ${binTarget} -> ${binSource}`);
        } catch (linkError) {
          console.log(`\n⚠️ 权限不足，无法自动创建 /usr/local/bin 软链。`);
          console.log(`请手动在终端执行以下命令（可能需要输入密码）:`);
          console.log(`\n  sudo ln -sf "${binSource}" "${binTarget}"\n`);
          
          // 如果这里没法创建，sokeclaw 可能还是找不到，提示一下
          console.log(`执行上述命令后，再重新打开 SokeClaw 即可。`);
          process.exit(0);
        }
      } else {
        console.log(`⚠️ 未找到源文件: ${binSource}`);
      }
      
      console.log('尝试重启 SokeClaw / WorkClaw...');
      try { execSync('killall SokeClaw 2>/dev/null'); } catch(e) {}
      try { execSync('killall WorkClaw 2>/dev/null'); } catch(e) {}
      try { execSync('killall Electron 2>/dev/null'); } catch(e) {}
      
      // 等待进程退出
      setTimeout(() => {
        let launched = false;
        if (fs.existsSync('/Applications/SokeClaw.app')) {
          execSync('open "/Applications/SokeClaw.app"');
          launched = true;
        } else if (fs.existsSync('/Applications/WorkClaw.app')) {
          execSync('open "/Applications/WorkClaw.app"');
          launched = true;
        }
        
        if (launched) {
          console.log('\n✅ 修复完成！SokeClaw 已重启。现在你应该可以在里面使用 /skill soke-exam 了。');
        } else {
          console.log('\n✅ 环境变量已注入，但未找到 SokeClaw 应用。请手动从“启动台(Launchpad)”或“访达(Finder)”重新打开它。');
        }
        process.exit(0);
      }, 1500);
      return; // 异步等待中
    } else {
      console.log('\n⚠️ 目前 setup-gui-env 仅支持 macOS 平台。Windows/Linux 用户请手动将 Node/npm bin 目录加入系统环境变量。');
      process.exit(0);
    }
  } catch (error) {
    console.error('修复过程中出现错误:', error.message);
    process.exit(1);
  }
}

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