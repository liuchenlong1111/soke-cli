# Soke CLI 简单开发流程

## 1. 准备开发环境
克隆/下载本仓库代码，并让 AI 协助你安装好 Go 环境。
```bash
git clone https://codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli.git
cd soke-cli

# 基于 master 分支创建你的个人开发分支
git checkout master
git fetch origin master
git merge origin/master
git branch my-cli-<你的名字>
git checkout -b my-cli-<你的名字>
```

## 2. 安装客户端与初始化配置
下载安装好 SokeClaw，并全局注册 Skill 供 Agent 调用：
```bash
npm install -g @sokeai/soke-claw@latest
npx skills add liuchenlong1111/soke-cli -y -g
```

**初始化 CLI 配置**
执行以下命令初始化权限：
```bash
soke-cli config init
```
按提示填写测试环境参数（暂时使用以下测试配置）：
- **app_key**: `soke426c576c4ce58c02e80f5f2de28b65bd`
- **app_secret**: `xxxxxxx` *(请向管理员获取)*
- **API地址**: `https://opendev.soke.cn`
- **corpid**: `dingc4cffb655c00d67f35c2f4657eb6378f`
- **dept_user_id**: `manager5600`

完成后，执行 `soke-cli config show` 检查配置是否正确。

## 3. 开发与本地测试
通过自然语言描述让 AI 协助你创建或修改 Skill。开发完成后，运行本地测试脚本，验证本地环境并将最新的 Skill 自动分发到本地 Agent：
```bash
./scripts/local-test.sh
```

## 4. 联调验证
重新打开 Agent（如 SokeClaw）进行实际对话测试。如果不满意，继续调整 Skill 并重复步骤 3，直到满意为止。

## 5. 提交代码与发布
开发并测试满意后，将代码推送到远程仓库：
```bash
git push origin my-cli-<你的名字>
```
最后，通知维护人员打包发版到 NPM 和 GitHub，然后分发到所有客户。