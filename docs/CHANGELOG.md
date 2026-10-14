# Changelog

All notable changes to GoTunnel are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and GoTunnel follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## How this file works

- Changes land under **[Unreleased]** as they are made, grouped by type: **Added**, **Changed**, **Fixed**, **Security**, **Removed**, **Breaking**.
- When a release is cut, the `[Unreleased]` heading is renamed to the new version with the release date, and a fresh empty `[Unreleased]` is started.
- Versions are ordered newest first. Each version links to its GitHub release.
- Version numbers follow SemVer: **MAJOR** for breaking changes, **MINOR** for backward-compatible features, **PATCH** for fixes.

---

## [Unreleased]

### Added

- **One-line install.** `curl -fsSL https://gotunnel.bakaredev.site/install.sh | sh` on macOS and Linux, `irm https://gotunnel.bakaredev.site/install.ps1 | iex` on Windows. The script picks the right build, verifies its SHA-256 checksum and puts `gotunnel` on your PATH. `GOTUNNEL_VERSION` pins a release and `GOTUNNEL_INSTALL_DIR` changes where it goes.
- **npm package.** `npm install -g gotunnel` (Node.js 22+), or `npx gotunnel http 3000` without installing. It ships the prebuilt binary for your platform and downloads nothing at install time.
- **`go install`.** `go install github.com/bakare-dev/gotunnel-io/cmd/gotunnel@latest` installs a small launcher that, on first run, downloads the matching release, verifies its checksum and replaces itself with it.
- **Download & install guide** in the docs, covering every install method, direct downloads, updating, uninstalling and troubleshooting.
- Each release now includes a `checksums.txt` with the SHA-256 of every binary.

### Changed

- `gotunnel update` and `gotunnel uninstall` recognise an npm install: `update` points you to `npm install -g gotunnel@latest`, and `uninstall` frees your registration and deletes your token but leaves removing the package to `npm uninstall -g gotunnel`.

---

## [1.0.0] - 2026-10-14

The first release of GoTunnel, a free personal side project by [BakareDev](https://bakaredev.site), provided as-is with fair-use limits on the hosted server.

### Tunnels

- **`gotunnel http <port>`** exposes a local web app at a public `https://<id>.gotunnel.bakaredev.site` URL. It takes a port or `host:port`, and a TLS backend with `--local https://localhost:PORT`.
- **`gotunnel tcp <port>`** exposes a database or any other TCP service (Postgres, MySQL, MongoDB, Redis, SSH). The person connecting uses their normal client, for example `psql -h <host> -p <port>`, with nothing to install.
- TCP tunnels are **locked by default**. Only IPs you pass with `--allowlist`, or that are approved through the printed approval link, can connect. The person connecting verifies their email with a one-time code, their current IP is approved until the tunnel ends, and you see each approval live in your terminal.
- Accounts get a stable name derived from their email, so a reconnect keeps the same address.
- When nothing is running on the local port, visitors see a GoTunnel page, `ERR_GOTUNNEL_LOCAL_UNAVAILABLE`, instead of a bare proxy error.

### Client

- First run registers the device: enter your Gmail address and the one-time code emailed to you. One device per address.
- **Request inspector**, on by default with `gotunnel http`: open `http://127.0.0.1:4040` to see every request with headers and body, **Replay** it against your local service, or **Copy as cURL**. History lives only in memory for that run. Turn it off with `--inspect=false`.
- Live terminal view with the endpoint, latency, uptime and a request log.
- Reconnects on its own after network drops and keeps trying until the server is back, backing off from 1 to 30 seconds with jitter. `--max-retries <n>` gives up after `n` attempts, and `--no-reconnect` turns reconnecting off.
- `gotunnel status`, `gotunnel audit` (PASS/WARN/FAIL security check of this machine), `gotunnel benchmark` (DNS, TCP, TLS and first-byte timing), `gotunnel inspect` and `gotunnel update`.
- `gotunnel uninstall` removes GoTunnel from the computer and frees your email on the hosted server, so you can register again on this or another device.
- Single binary for macOS, Linux and Windows.

### Hosted server

- Free with fair-use limits: 30-minute sessions, a daily session limit, one HTTP tunnel per device and one TCP tunnel per account.
- Clients connect over WebSocket on port 443, so it works from restrictive networks.

### Self-hosting

- `gotunnel server` runs your own server in port mode or subdomain mode, with static tokens (`--auth-tokens`) or email registration.
- TLS on the control channel, optional mutual TLS, per-tunnel bandwidth limits, a concurrent-tunnel cap, a daily session limit with exempt emails, and a YAML config file with flags > environment > file precedence.
- `/healthz` and Prometheus `/metrics` for tunnels, streams, auth results and rejections.
- TCP tunnels with approval links, an audit log of TCP sessions, approvals and connections (metadata only, never traffic), and optional nftables gating so the TCP port range is closed except to approved IPs.
- `gotunnel admin` commands over a local owner-only socket: `sessions`, `kill`, `ban`, `unban`, `bans`. Abuse alert emails include ready-to-run kill and ban commands.

### Security

- Tokens are HMAC-signed and bound to the device, compared in constant time and never logged. Repeated failed logins lock the IP out.
- Inbound frames are checked against the handshake → auth → forwarding state machine, with frame size limits and a per-tunnel stream cap.
- Read and write deadlines on public and local connections.
- IPs that keep probing TCP tunnels without approval are temporarily banned. Banned accounts are refused when connecting and when registering.
- The request inspector escapes captured data, only answers localhost, and refuses cross-site writes.

---

[Unreleased]: https://github.com/bakare-dev/gotunnel-io/releases
[1.0.0]: https://github.com/bakare-dev/gotunnel-io/releases/tag/v1.0.0
