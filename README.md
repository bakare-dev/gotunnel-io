# GoTunnel

[![Download GoTunnel](https://img.shields.io/badge/Download-GoTunnel-0A66C2?style=for-the-badge)](https://gotunnel.bakaredev.site/download)
[![Latest Release](https://img.shields.io/github/v/release/bakare-dev/gotunnel-io?style=for-the-badge&label=Latest%20Release)](https://github.com/bakare-dev/gotunnel-io/releases/latest)
[![Personal Project](https://img.shields.io/badge/Status-Personal%20Project-22C55E?style=for-the-badge)](https://gotunnel.bakaredev.site)

HTTP and TCP tunnels that put a local app or database on the internet with one command.

This repository hosts the **prebuilt GoTunnel binaries** and docs. Use the download page for the recommended installer for your platform, or grab a release artifact directly from GitHub.

> **A side project, not a company.** GoTunnel is a personal project by [BakareDev](https://bakaredev.site), built to learn networking and systems design in Go. The hosted server is free with fair-use limits (30-minute sessions, 3 per day, one tunnel per device), provided as-is without guarantees, and may change or go offline. For anything important, [self-host it](docs/SELF_HOSTING.md).

## Download

Get the latest GoTunnel build from the official download page:

**[Download GoTunnel](https://gotunnel.bakaredev.site/download)**

Direct release artifacts are also available below:

| Platform | Architecture | Download |
| --- | --- | --- |
| macOS | Apple Silicon (M1+) | [gotunnel-darwin-arm64](https://github.com/bakare-dev/gotunnel-io/releases/latest/download/gotunnel-darwin-arm64) |
| macOS | Intel | [gotunnel-darwin-amd64](https://github.com/bakare-dev/gotunnel-io/releases/latest/download/gotunnel-darwin-amd64) |
| Linux | x86-64 | [gotunnel-linux-amd64](https://github.com/bakare-dev/gotunnel-io/releases/latest/download/gotunnel-linux-amd64) |
| Linux | ARM64 | [gotunnel-linux-arm64](https://github.com/bakare-dev/gotunnel-io/releases/latest/download/gotunnel-linux-arm64) |
| Windows | x86-64 | [gotunnel-windows-amd64.exe](https://github.com/bakare-dev/gotunnel-io/releases/latest/download/gotunnel-windows-amd64.exe) |

Or see all builds on the [releases page](https://github.com/bakare-dev/gotunnel-io/releases/latest).

## Install

### macOS (Apple Silicon)

```bash
curl -L -o gotunnel https://github.com/bakare-dev/gotunnel-io/releases/latest/download/gotunnel-darwin-arm64
chmod +x gotunnel
sudo mv gotunnel /usr/local/bin/gotunnel
xattr -dr com.apple.quarantine /usr/local/bin/gotunnel   # if Gatekeeper blocks it
```

(For Intel Macs, swap `gotunnel-darwin-arm64` for `gotunnel-darwin-amd64`.)

### Linux (x86-64)

```bash
curl -L -o gotunnel https://github.com/bakare-dev/gotunnel-io/releases/latest/download/gotunnel-linux-amd64
chmod +x gotunnel
sudo mv gotunnel /usr/local/bin/gotunnel
```

(For ARM64, use `gotunnel-linux-arm64`.)

### Windows

Download [gotunnel-windows-amd64.exe](https://github.com/bakare-dev/gotunnel-io/releases/latest/download/gotunnel-windows-amd64.exe), then in PowerShell install it to a per-user folder and add that folder to your `PATH`:

```powershell
$dir = "$env:LOCALAPPDATA\Programs\gotunnel"
New-Item -ItemType Directory -Force -Path $dir | Out-Null
Move-Item -Force "$HOME\Downloads\gotunnel-windows-amd64.exe" "$dir\gotunnel.exe"
$old = [Environment]::GetEnvironmentVariable("Path","User")
if ($old -notlike "*$dir*") { [Environment]::SetEnvironmentVariable("Path", "$old;$dir", "User") }
$env:Path += ";$dir"
gotunnel version
```

The `$env:Path += ";$dir"` line makes `gotunnel` work in the current window; the `SetEnvironmentVariable` line makes it stick for new terminals. If `gotunnel` still isn't found, open a new PowerShell window.

## Quick start

Expose a local service (e.g. running on port 3000) through the free hosted server:

```bash
gotunnel http 3000
```

The first time you run it on a device, GoTunnel asks for your Gmail address and emails you a one-time code to register the device (one device per address). Your email is only used for that. After that, the device stays signed in and connects straight away.

You'll get a public `https://<random>.gotunnel.bakaredev.site` URL forwarding to your local port. Press `q` or `Ctrl-C` to stop.

### Share a database or other TCP service

```bash
gotunnel tcp 5432
```

You get a host and port, e.g. `k3x9m2qz7a.tcp.gotunnel.bakaredev.site:20017`, that works with any normal client: `psql`, `mysql`, `mongosh`, `redis-cli`, `ssh`. TCP tunnels are **locked**: only IPs you allow with `--allowlist`, or that are approved through the link GoTunnel prints (the person connecting verifies their email once), can connect. See [`gotunnel tcp`](docs/CLI.md#gotunnel-tcp-port).

## Update

Upgrade to the latest release in place:

```bash
gotunnel update
```

## Uninstall

Remove GoTunnel from your computer:

```bash
gotunnel uninstall
```

It frees your email on the hosted server so you can register again later, on this or another device. Then it deletes your token and device ID, and the `gotunnel` binary. On macOS and Linux it asks for your password if the binary is in `/usr/local/bin`. On Windows it also removes the install folder from your `PATH`. See [`gotunnel uninstall`](docs/CLI.md#gotunnel-uninstall) for the options.

## Commands

```
gotunnel http <port>                                Start a tunnel on the hosted server
gotunnel tcp <port> [--allowlist ip,cidr]           Share a database or other TCP service, locked to approved IPs
gotunnel --server <host> --local localhost:<port>   Start a tunnel against any server
gotunnel status                                     Registration, token expiry, server reachability
gotunnel audit                                      Security check of this machine's setup
gotunnel benchmark                                  DNS / TCP / TLS / first-byte latency to the server
gotunnel inspect                                    What this build supports and where its files live
gotunnel update                                     Update to the latest release
gotunnel uninstall                                  Remove GoTunnel and free your email for re-registration
gotunnel version                                    Show the installed version
gotunnel help                                       Show all options
```

## Documentation

The docs below explain how GoTunnel works so you can understand what you're running before you trust it with your traffic.

| Document | What it covers |
| --- | --- |
| [Command reference](docs/CLI.md) | Every command: `http`, `tcp`, `status`, `audit`, `benchmark`, `inspect`, `update`, `uninstall`, `admin` |
| [Self-Hosting](docs/SELF_HOSTING.md) | Run your own GoTunnel server on a VPS instead of the hosted server |
| [Architecture](docs/ARCHITECTURE.md) | How the client, server, tunnels, and security model fit together |
| [Protocol](docs/PROTOCOL.md) | The binary wire protocol: frames, message types, and the session state machine |
| [Changelog](docs/CHANGELOG.md) | What changed in each release |

## Feedback and issues

Found a bug or have a request? [Open an issue](https://github.com/bakare-dev/gotunnel-io/issues/new/choose). If GoTunnel is useful to you, consider [starring the repo](https://github.com/bakare-dev/gotunnel-io) to follow releases.
