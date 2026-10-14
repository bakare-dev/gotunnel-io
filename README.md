# GoTunnel

[![Download GoTunnel](https://img.shields.io/badge/Download-GoTunnel-0A66C2?style=for-the-badge)](https://gotunnel.bakaredev.site/download)
[![Latest Release](https://img.shields.io/github/v/release/bakare-dev/gotunnel-io?style=for-the-badge&label=Latest%20Release)](https://github.com/bakare-dev/gotunnel-io/releases/latest)
[![Personal Project](https://img.shields.io/badge/Status-Personal%20Project-22C55E?style=for-the-badge)](https://gotunnel.bakaredev.site)

HTTP and TCP tunnels that put a local app or database on the internet with one command.

This repository hosts the **prebuilt GoTunnel binaries**, the docs, and the small launcher behind `go install`.

> **A side project, not a company.** GoTunnel is a personal project by [BakareDev](https://bakaredev.site), built to learn networking and systems design in Go. The hosted server is free with fair-use limits (30-minute sessions, 3 per day, one tunnel per device), provided as-is without guarantees, and may change or go offline. For anything important, [self-host it](docs/SELF_HOSTING.md).

## Install

macOS and Linux:

```bash
curl -fsSL https://gotunnel.bakaredev.site/install.sh | sh
```

Windows (PowerShell):

```powershell
irm https://gotunnel.bakaredev.site/install.ps1 | iex
```

With npm (Node.js 22+):

```bash
npm install -g gotunnel
```

With Go (1.21+):

```bash
go install github.com/bakare-dev/gotunnel-io/cmd/gotunnel@latest
```

Pick one method; you only need GoTunnel installed once.

Every method checks the download against the release's SHA-256 `checksums.txt`. The **[install guide](docs/INSTALL.md)** covers each one step by step, plus direct binary downloads, manual install, updating, uninstalling and troubleshooting. All builds are also on the [releases page](https://github.com/bakare-dev/gotunnel-io/releases/latest).

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

If you installed with npm, run `npm install -g gotunnel@latest` instead. See [Update](docs/INSTALL.md#update).

## Uninstall

Remove GoTunnel from your computer:

```bash
gotunnel uninstall
```

It frees your email on the hosted server so you can register again later, on this or another device. Then it deletes your token and device ID, and the `gotunnel` binary. On macOS and Linux it asks for your password if the binary is in `/usr/local/bin`. On Windows it also removes the install folder from your `PATH`. If you installed with npm, it leaves the package alone; finish with `npm uninstall -g gotunnel`. See [`gotunnel uninstall`](docs/CLI.md#gotunnel-uninstall) for the options.

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
| [Download & install](docs/INSTALL.md) | Every install method, direct downloads, updating, uninstalling and troubleshooting |
| [Command reference](docs/CLI.md) | Every command: `http`, `tcp`, `status`, `audit`, `benchmark`, `inspect`, `update`, `uninstall`, `admin` |
| [Self-Hosting](docs/SELF_HOSTING.md) | Run your own GoTunnel server on a VPS instead of the hosted server |
| [Architecture](docs/ARCHITECTURE.md) | How the client, server, tunnels, and security model fit together |
| [Protocol](docs/PROTOCOL.md) | The binary wire protocol: frames, message types, and the session state machine |
| [Changelog](docs/CHANGELOG.md) | What changed in each release |

## Feedback and issues

Found a bug or have a request? [Open an issue](https://github.com/bakare-dev/gotunnel-io/issues/new/choose). If GoTunnel is useful to you, consider [starring the repo](https://github.com/bakare-dev/gotunnel-io) to follow releases.
