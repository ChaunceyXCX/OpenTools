# ZTools (Wails Port)

<div align="center">

<img src="./.github/assets/icon.png" alt="ZTools Logo" width="120">

**一个高性能、可扩展的应用启动器和插件平台**

_uTools 的开源实现 | 支持 macOS、Windows 和 Linux_

[![License](https://img.shields.io/github/license/lzx8589561/ZTools)](./LICENSE)
[![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Windows%20%7C%20Linux-blue)](https://github.com/ChaunceyXCX/OpenTools)

[English](./README_EN.md) | 简体中文

</div>

---

## ✨ 特性

- 🚀 **快速启动** - 拼音搜索、正则匹配、历史记录、固定应用
- 🧩 **插件系统** - 支持 JS 脚本插件和 Go Service 插件，完整 API
- 📋 **剪贴板管理** - 历史记录、搜索、自动监控
- 🎨 **主题定制** - CSS 变量驱动的暗色主题
- ⚡ **高性能** - bbolt 数据库、Go 原生性能、极小内存占用
- 🌍 **跨平台** - 原生支持 macOS、Windows 和 Linux
- 🔒 **数据隔离** - 插件数据独立存储，安全可靠
- 🛠️ **开发友好** - Hot Reload 前端、Go Binding 类型安全
- ⚙️ **技术栈** - Wails v3 + Go 1.25 + Vue 3 + TypeScript
- 🧩 **浏览器自动化** - ZBrowser（chromedp）
- 🔌 **MCP Server** - AI 工具集成（scan/launch/clipboard）
- ☁ **WebDAV 同步** - 多端数据同步

## 🚀 快速开始

### 安装

> 本项目是 ZTools 的 **Wails 移植版**（`wails` 分支），使用 Go + WebView 替代了原 Electron 架构。
> 原始 Electron 版本在 `main` 分支。

#### Linux 依赖

```bash
# Ubuntu/Debian
sudo apt install libgtk-3-dev libwebkit2gtk-4.1-dev libx11-dev

# Fedora
sudo dnf install gtk3-devel webkit2gtk4.1-devel libX11-devel

# Arch
sudo pacman -S gtk3 webkit2gtk-4.1 libx11
```

#### 从源码构建

```bash
# 克隆仓库（wails 分支）
git clone -b wails https://github.com/ChaunceyXCX/OpenTools.git
cd OpenTools

# 安装前端依赖
cd frontend && pnpm install && cd ..

# 开发模式
wails3 dev

# 构建
wails3 build                    # 当前平台构建
wails3 build -platform linux    # Linux
wails3 build -platform darwin   # macOS
wails3 build -platform windows  # Windows
```

### 使用

1. 启动应用后，使用 `Alt+Z` 唤出主界面
2. 输入应用名称或命令搜索
3. `↑` `↓` 选择，`Enter` 启动，`Esc` 隐藏
4. 支持 Ctrl+Tab 在 Search / Clipboard / Settings 间切换

### 快捷键

| 快捷键 | 功能 |
|--------|------|
| `Alt+Z` | 唤出/隐藏主窗口 |
| `Esc` | 隐藏窗口 |
| `↑` `↓` | 导航搜索结果 |
| `Enter` | 启动选中应用 |

## 🧩 插件系统

ZTools 的 Wails 移植版提供与上游完全兼容的插件 API：

### 插件类型

- **Go Service 插件** - 编译时注册，高性能（如 system 插件）
- **JS 脚本插件** - goja 运行时执行，支持 `window.exports[code].args.enter()` 模式
- **ZPX 包插件** - 完全兼容上游 `.zpx` 格式（asar+gzip）

### 插件 API

插件可通过 `window.ztools.*` 调用以下能力：

| API | 说明 |
|-----|------|
| `ztools.db.*` | 数据库 CRUD（同步 + Promise 异步） |
| `ztools.clipboard.*` | 剪贴板读写 |
| `ztools.shell.*` | Shell 命令执行 |
| `ztools.showNotification()` | 系统通知 |
| `ztools.showToast()` | Toast 提示 |
| `ztools.setSubInput()` | 子输入框（插件搜索） |
| `ztools.hideWindow()` | 隐藏窗口 |
| `ztools.simulateKeyboardTap()` | 键盘模拟 |

## 🗂 项目结构

```
OpenTools/
├── main.go                  # Go 应用入口
├── internal/                # Go 后端
│   ├── api/                 # Wails Binding 服务
│   ├── core/                # 核心模块
│   │   ├── database/        # bbolt 数据库
│   │   ├── clipboard/       # 剪贴板监控
│   │   ├── scanner/         # 系统应用扫描
│   │   ├── launcher/        # 指令启动
│   │   ├── mcp/             # MCP Server
│   │   ├── httpserver/      # HTTP API
│   │   ├── zbrowser/        # 浏览器自动化
│   │   ├── sync/            # WebDAV 同步
│   │   └── translation/     # 离线翻译
│   ├── plugin/              # 插件运行时
│   └── native/              # 原生能力
├── frontend/                # Vue 3 前端
│   ├── src/
│   │   ├── components/      # UI 组件
│   │   ├── stores/          # Pinia 状态管理
│   │   └── style.css        # 全局样式
│   └── bindings/            # Wails 自动生成绑定
└── build/                   # 打包配置
```

## 📦 技术栈

| 组件 | 技术 |
|------|------|
| 桌面框架 | Wails v3 (alpha.90) |
| 后端语言 | Go 1.25 |
| 前端 | Vue 3 + TypeScript |
| 状态管理 | Pinia |
| 搜索 | Fuse.js |
| 数据库 | bbolt (LMDB 兼容) |
| JS 引擎 | goja（插件运行时） |
| 浏览器自动化 | chromedp |
| 剪贴板 | atotto/clipboard |
| WebDAV | gowebdav |
| 打包 | Wails 原生构建 |

## 🔄 分支策略

| 分支 | 说明 |
|------|------|
| `main` | 上游同步分支，与原 ZTools 一致 |
| `wails` | Wails 重构分支（当前） |

`wails` 分支定期从 `main` merge 上游更新，新功能同步移植到 Go 实现。

> **关于项目名称变更**：
> 本项目从 [ZToolsCenter/ZTools](https://github.com/ZToolsCenter/ZTools) fork 而来。`main` 分支保持与原上游同步，
> `wails` 分支在重构的同时将应用名改为 **OpenTools**，以标识这是一个独立的 Wails 移植版本。
>
> 重命名涉及范围：应用名、窗口标题、构建产物、数据目录 (`~/.ztools/` → `~/.opentools/`)、
> 产品标识符 (`link.eiot.ztools` → `com.opentools.app`)、API 返回值、Linux 桌面文件、MCP 服务名等。
>
> **保留不变**：`window.ztools.*` 插件 API 名称（插件生态标准）、`ZTOOLS/` 数据库命名空间（数据兼容）。

## 📄 许可证

[MIT](./LICENSE)
