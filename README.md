<h1 align="center">Usque Android App / Usque 安卓应用</h1>

<p align="center">
  <strong>Cloudflare WARP / MASQUE VPN Client for Android</strong><br>
  <strong>Android 版 Cloudflare WARP / MASQUE VPN 客户端</strong>
</p>

<p align="center">
  <a href="https://github.com/garthnet/usque-android-app/releases">
    <img src="https://img.shields.io/github/v/release/garthnet/usque-android-app?style=flat-square" alt="Release">
  </a>
  <a href="https://github.com/garthnet/usque-android-app">
    <img src="https://img.shields.io/badge/platform-android-blue?style=flat-square" alt="Platform">
  </a>
  <a href="https://github.com/garthnet/usque-android-app/actions/workflows/android.yml">
    <img src="https://img.shields.io/github/actions/workflow/status/garthnet/usque-android-app/android.yml?branch=main&style=flat-square" alt="Android APK">
  </a>
</p>

<p align="center">
  <a href="#english">English</a> · <a href="#中文">中文</a>
</p>

---

## English

Usque Android App is a native Android VPN app for Cloudflare WARP / Zero Trust. It wraps the MASQUE-based Go client from the upstream [`usque`](https://github.com/Diniboy1123/usque) project and provides an Android UI, `VpnService` integration, per-app routing, and automated APK releases.

- **Android package**: `com.warp.usque`
- **Latest release**: <https://github.com/garthnet/usque-android-app/releases/tag/v1.0.2>
- **Latest APK**: <https://github.com/garthnet/usque-android-app/releases/download/v1.0.2/usque-android-app-release-v1.0.2.apk>

## ✨ Features

- 🚀 **Cloudflare WARP** - MASQUE-based WARP tunnel using the upstream Go core
- 🔒 **Global VPN** - Route all app traffic through the VPN
- 📱 **Per-App VPN** - Route only selected apps through the VPN
- 🌍 **Bilingual UI** - Chinese / English interface toggle
- ⚙️ **Profiles** - Save and apply multiple endpoint / SNI profiles
- 📊 **Speed Display** - Runtime traffic speed using Android `TrafficStats`
- 🛡️ **Foreground Service** - Better long-running VPN stability on Android
- 🔁 **Auto Restart** - Automatically restarts the native tunnel after unexpected disconnects
- ⚡ **Native Core** - Go + gomobile AAR integration for native performance

## 📲 Download

Download the latest APK from the [Releases](https://github.com/garthnet/usque-android-app/releases) page.

Direct download:

```text
https://github.com/garthnet/usque-android-app/releases/download/v1.0.2/usque-android-app-release-v1.0.2.apk
```

> The GitHub Actions workflow reads `versionName` from `app/build.gradle` and publishes the APK to the matching versioned Release, for example `1.0.2` → `v1.0.2`.

## 🛠️ Build from Source

### Prerequisites

- JDK 17
- Android SDK, API 34
- Gradle 8.5 or compatible Gradle setup
- Prebuilt upstream Android artifacts already included in this repository:
  - `app/libs/usque.aar`
  - `app/libs/usque-classes.jar`

### Build Steps

```bash
# Debug APK
gradle clean assembleDebug

# Release APK
gradle clean assembleRelease
```

If no release signing config is provided, the release APK may be unsigned depending on the environment.

### GitHub Actions Recommended

Push to `main` or manually run the Android APK workflow:

```bash
gh workflow run android.yml
```

The workflow will:

- Build with JDK 17 and Gradle 8.5
- Run `gradle clean` to avoid stale dex / manifest output
- Build a signed release APK when signing secrets are configured
- Upload the APK and `SHA256SUMS` to the matching GitHub Release

Required signing secrets:

```text
ANDROID_KEYSTORE_BASE64
ANDROID_KEYSTORE_PASSWORD
ANDROID_KEY_ALIAS
ANDROID_KEY_PASSWORD
```

See also:

```text
docs/github-actions-signing.md
```

## 📖 Usage

1. **Install APK** - Download and install the latest APK on an Android device.
2. **First Run** - The app registers a Cloudflare WARP account automatically if no valid config exists.
3. **Configure** - Set SNI, endpoint, port, and optional profiles.
4. **Choose Mode** - Use Global mode or Per-App mode.
5. **Connect** - Tap the connect button to start the VPN.

### Proxy Modes

| Mode | Description |
|------|-------------|
| **Global** | All app traffic goes through the VPN. The app itself is excluded to avoid routing the control connection back into its own tunnel. |
| **Per-App** | Only selected apps use the VPN. Other apps connect directly. The app itself is skipped for the same control-connection reason. |

### Settings

| Parameter | Example / Default | Description |
|-----------|-------------------|-------------|
| SNI | `www.visa.cn` | TLS SNI value used by the tunnel |
| Endpoint | `162.159.198.2:500` | WARP / MASQUE endpoint |
| Mode | Global / Per-App | Traffic routing mode |

## 🏗️ Architecture

```text
usque-android-app/
├── app/
│   ├── src/main/
│   │   ├── kotlin/com/warp/usque/      # Kotlin UI and VPN service
│   │   ├── res/                        # Android resources
│   │   └── AndroidManifest.xml
│   ├── libs/                           # Prebuilt upstream AAR / Java classes
│   │   ├── usque.aar
│   │   └── usque-classes.jar
│   └── build.gradle
├── docs/
│   └── github-actions-signing.md
└── .github/workflows/
    └── android.yml                     # CI APK build and Release upload
```

### Tech Stack

- **UI**: Kotlin + Android Views
- **VPN Service**: Android `VpnService` + real TUN fd
- **VPN Core**: Go + gomobile → AAR
- **Protocol**: Cloudflare MASQUE / WARP
- **Build**: Gradle + GitHub Actions

## 🎨 Design

- **Style**: Material-style Android UI
- **Theme**: Soft warm orange / light card layout
- **Modes**: Overview, Config, Apps
- **UX**: Mobile-first controls, searchable app list, bilingual labels

## 📝 Notes

- The app uses `VpnService.Builder.establish()` to obtain a real Android TUN fd.
- The detached fd is passed to the upstream native layer through `startTunnelWithFd`.
- Endpoint and SNI are configured through the `usqueandroid` API before the tunnel starts.
- The VPN service runs as a foreground service for better long-running stability.
- Unexpected native tunnel disconnects trigger automatic restart unless the user manually stopped the VPN.
- GitHub Actions uses `gradle clean` before building to prevent stale dex / manifest mismatches after package-name changes.

## 📝 License

This project is based on [`usque`](https://github.com/Diniboy1123/usque) by Diniboy1123. Check the upstream project license and comply with its terms when redistributing the combined work.

## 🙏 Acknowledgements

- [`usque`](https://github.com/Diniboy1123/usque) - Original WARP / MASQUE implementation
- [Cloudflare WARP](https://developers.cloudflare.com/warp-client/) - WARP platform
- [MASQUE](https://datatracker.ietf.org/doc/rfc9484/) - RFC 9484

---

## 中文

Usque 安卓应用是一个面向 Cloudflare WARP / Zero Trust 的原生 Android VPN 应用。它基于上游 [`usque`](https://github.com/Diniboy1123/usque) 项目的 MASQUE Go 客户端，提供 Android 图形界面、`VpnService` 集成、分应用路由和 APK 自动发布流程。

- **Android 包名**：`com.warp.usque`
- **最新版本**：<https://github.com/garthnet/usque-android-app/releases/tag/v1.0.2>
- **最新 APK**：<https://github.com/garthnet/usque-android-app/releases/download/v1.0.2/usque-android-app-release-v1.0.2.apk>

## ✨ 功能特性

- 🚀 **Cloudflare WARP** - 基于上游 Go 核心的 MASQUE WARP 隧道
- 🔒 **全局 VPN** - 所有应用流量通过 VPN
- 📱 **分应用 VPN** - 仅选中的应用走 VPN，其它应用直连
- 🌍 **中英切换** - 支持中文 / 英文界面
- ⚙️ **配置档** - 支持保存和应用多个 endpoint / SNI 配置
- 📊 **速度显示** - 基于 Android `TrafficStats` 显示实时速度
- 🛡️ **前台服务** - 提升 Android 长时间后台 VPN 稳定性
- 🔁 **自动重连** - native 隧道异常断开后自动重启连接
- ⚡ **原生核心** - Go + gomobile AAR 集成，接近原生性能

## 📲 下载

从 [Releases](https://github.com/garthnet/usque-android-app/releases) 页面下载最新版 APK。

直接下载：

```text
https://github.com/garthnet/usque-android-app/releases/download/v1.0.2/usque-android-app-release-v1.0.2.apk
```

> GitHub Actions 会从 `app/build.gradle` 读取 `versionName`，并把 APK 发布到对应版本 Release，例如 `1.0.2` → `v1.0.2`。

## 🛠️ 从源码构建

### 环境要求

- JDK 17
- Android SDK，API 34
- Gradle 8.5 或兼容的 Gradle 环境
- 本仓库已包含上游预构建 Android 产物：
  - `app/libs/usque.aar`
  - `app/libs/usque-classes.jar`

### 构建步骤

```bash
# Debug APK
gradle clean assembleDebug

# Release APK
gradle clean assembleRelease
```

如果没有配置 release 签名，本地 release APK 可能会是未签名产物，具体取决于构建环境。

### 推荐使用 GitHub Actions

推送到 `main`，或手动运行 Android APK 工作流：

```bash
gh workflow run android.yml
```

工作流会：

- 使用 JDK 17 和 Gradle 8.5 构建
- 执行 `gradle clean`，避免旧 dex / 新 Manifest 混用
- 配置签名 secrets 时构建签名 release APK
- 上传 APK 和 `SHA256SUMS` 到对应 GitHub Release

签名 release 构建需要配置：

```text
ANDROID_KEYSTORE_BASE64
ANDROID_KEYSTORE_PASSWORD
ANDROID_KEY_ALIAS
ANDROID_KEY_PASSWORD
```

更多说明：

```text
docs/github-actions-signing.md
```

## 📖 使用方法

1. **安装 APK** - 下载最新版 APK 并安装到 Android 设备。
2. **首次运行** - 如果没有有效配置，应用会自动注册 Cloudflare WARP 账户。
3. **配置参数** - 设置 SNI、endpoint、端口，也可以保存配置档。
4. **选择模式** - 使用全局模式或分应用模式。
5. **连接 VPN** - 点击连接按钮开始使用。

### 代理模式

| 模式 | 说明 |
|------|------|
| **Global / 全局** | 所有应用流量走 VPN。应用自身会被排除，避免控制连接被套回自己的 VPN 隧道。 |
| **Per-App / 分应用** | 只有选中的应用走 VPN，其它应用直连。应用自身也会被跳过，原因相同。 |

### 设置项

| 参数 | 示例 / 默认值 | 说明 |
|------|---------------|------|
| SNI | `www.visa.cn` | 隧道使用的 TLS SNI 值 |
| Endpoint | `162.159.198.2:500` | WARP / MASQUE endpoint |
| Mode | Global / Per-App | 流量路由模式 |

## 🏗️ 架构

```text
usque-android-app/
├── app/
│   ├── src/main/
│   │   ├── kotlin/com/warp/usque/      # Kotlin UI 和 VPN 服务
│   │   ├── res/                        # Android 资源
│   │   └── AndroidManifest.xml
│   ├── libs/                           # 上游预构建 AAR / Java 类
│   │   ├── usque.aar
│   │   └── usque-classes.jar
│   └── build.gradle
├── docs/
│   └── github-actions-signing.md
└── .github/workflows/
    └── android.yml                     # CI APK 构建与 Release 上传
```

### 技术栈

- **UI**：Kotlin + Android Views
- **VPN 服务**：Android `VpnService` + 真实 TUN fd
- **VPN 核心**：Go + gomobile → AAR
- **协议**：Cloudflare MASQUE / WARP
- **构建**：Gradle + GitHub Actions

## 🎨 设计

- **风格**：Material-style Android UI
- **主题**：柔和暖橙色 / 浅色卡片布局
- **页面**：总览、配置、应用
- **体验**：移动端优先、应用列表搜索、中英双语标签

## 📝 说明

- 应用通过 `VpnService.Builder.establish()` 获取真实 Android TUN fd。
- detached fd 会通过 `startTunnelWithFd` 传给上游 native 层。
- 启动隧道前，会通过 `usqueandroid` API 设置 endpoint 和 SNI。
- VPN 服务以前台服务运行，提升长时间后台稳定性。
- 除非用户手动停止，native 隧道异常断开后会自动重启。
- GitHub Actions 构建前执行 `gradle clean`，避免修改包名后出现旧 dex / 新 Manifest 混用。

## 📝 许可证

本项目基于 Diniboy1123 的 [`usque`](https://github.com/Diniboy1123/usque) 项目封装。重新分发组合产物时，请检查并遵守上游项目许可证。

## 🙏 致谢

- [`usque`](https://github.com/Diniboy1123/usque) - 原始 WARP / MASQUE 实现
- [Cloudflare WARP](https://developers.cloudflare.com/warp-client/) - WARP 平台
- [MASQUE](https://datatracker.ietf.org/doc/rfc9484/) - RFC 9484

---
