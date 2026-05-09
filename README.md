# Usque Android App

A native Android VPN app for Cloudflare WARP / Zero Trust, built around the MASQUE-based Go client from:

<https://github.com/Diniboy1123/usque>

This repository contains the Android wrapper, UI, `VpnService` integration, and the prebuilt Android library artifacts required to build the APK.

## Upstream attribution

The core WARP / MASQUE client logic comes from the upstream `usque` project:

- Upstream repository: <https://github.com/Diniboy1123/usque>
- Upstream language: Go
- Android integration: prebuilt `usque.aar` / `usque-classes.jar` exposed through the `usqueandroid` package

This Android app adds:

- Kotlin Android UI
- Android `VpnService` lifecycle management
- TUN fd handoff through `startTunnelWithFd`
- Cloudflare WARP registration and connection controls
- Global proxy / per-app proxy mode
- Multi-profile connection settings
- Mobile-first Material-style UI

## Features

- Native Android VPN implementation using `VpnService`
- Cloudflare WARP / MASQUE connection through the upstream `usque` core
- Automatic first-run registration
- Global proxy and selected-app proxy modes
- App selection list with search, select all, and clear selection
- Multiple connection profiles
- Chinese / English language toggle
- Runtime speed display based on Android `TrafficStats`

## Repository contents

Important tracked files:

```text
app/src/main/kotlin/com/diniboy/usqueandroid/MainActivity.kt
app/src/main/kotlin/com/diniboy/usqueandroid/UsqueVpnService.kt
app/libs/usque.aar
app/libs/usque-classes.jar
.github/workflows/android.yml
docs/github-actions-signing.md
```

The repository intentionally does **not** include stale Speedgo test code, local probes, APK build outputs, local Gradle caches, `local.properties`, or signing keystores.

## Building locally

Use JDK 17 and Gradle / Android Gradle Plugin compatible with this project.

```bash
gradle assembleDebug
```

Release build:

```bash
gradle assembleRelease
```

If no production signing config is provided, release APK output may be unsigned depending on the environment.

## GitHub Actions APK build

The repository includes a GitHub Actions workflow at:

```text
.github/workflows/android.yml
```

Behavior:

- Without signing secrets: builds a debug APK artifact.
- With signing secrets: builds a signed release APK artifact.
- On GitHub Release publish: uploads APK artifacts to the Release.

See signing setup instructions:

```text
docs/github-actions-signing.md
```

Do **not** commit keystore files to this repository. Use GitHub Actions Secrets instead.

## Signing secrets

For signed release builds in GitHub Actions, configure these repository secrets:

```text
ANDROID_KEYSTORE_BASE64
ANDROID_KEYSTORE_PASSWORD
ANDROID_KEY_ALIAS
ANDROID_KEY_PASSWORD
```

## Architecture

```text
Android App / Kotlin UI
        │
        ▼
Android VpnService + real TUN fd
        │
        ▼
usqueandroid package from usque.aar / usque-classes.jar
        │
        ▼
Upstream Go usque core
        │
        ▼
Cloudflare WARP / MASQUE
```

## Notes

- The app uses `builder.establish()` to obtain a real Android TUN fd.
- The detached fd is passed to the upstream native layer via `startTunnelWithFd`.
- Endpoint is configured separately through the `usqueandroid` API before starting the tunnel.
- Per-app proxy mode uses Android `VpnService.Builder.addAllowedApplication()`.

## License

This repository is an Android app wrapper around the upstream `usque` project. Check the upstream project license and comply with its terms when redistributing the combined work:

<https://github.com/Diniboy1123/usque>
