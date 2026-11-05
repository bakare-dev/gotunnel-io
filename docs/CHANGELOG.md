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

---

## [1.3.0] - 2026-11-05

### Security

- **HTTP tunnels:** the first time someone opens an HTTP tunnel in a browser, they see a warning page explaining that it is served from someone's own computer. They click **Continue to site** to go through, and that choice is remembered for 7 days on that tunnel. API clients, webhooks, scripts and in-page requests are not affected. Browser tools can skip the page by sending the `Gotunnel-Skip-Browser-Warning` header.
- **Self-hosting:** turn the warning page on with `GOTUNNEL_BROWSER_WARNING=true`. It is off by default.
- **Saved tokens stay with their server.** The token saved at registration is now sent only to the server that issued it, and only over an encrypted connection. `gotunnel uninstall` follows the same rule.
- **Released and replaced tokens stop working.** After `gotunnel uninstall`, or after registering the same address again, older tokens for that device are refused at connect.
- **Registration limits.** Verification emails and code attempts are now limited per address and per client, whichever device ID is used. Registration responses no longer reveal whether an address is already registered or suspended; the explanation goes to that inbox instead.
- **Tunnel addresses no longer reveal the owner.** Stable subdomains are derived with a server-side secret, so nobody can work out your tunnel URL from your email.
- **Warning page.** Form submissions from another site now see the warning page too, and your choice to continue is tied to your IP address, so you may see the page again when your network changes. Sign-in callbacks from identity providers still work once you have continued.
- **Connection limits.** The server now disconnects connections that don't authenticate within 15 seconds, limits unauthenticated connections per IP, and refuses oversized frames before authentication.
- **Failed account checks now refuse the connection.** If the server can't confirm an account isn't banned, it asks the client to try again shortly instead of letting it connect.
- Visitor-controlled text in request logs and server messages is now escaped, so it can't forge log lines or send control codes to your terminal.
- **Request inspector.** The inspector now needs the link GoTunnel prints, with a token that changes on every run, so other users on the same machine can't read your traffic. It can no longer be framed by another page, and **Copy as cURL** leaves out unusual request methods so a pasted command can only run curl.
- **Control channel.** Browsers can no longer open the WebSocket control channel.
- **Signed releases.** Each release's `checksums.txt` is now signed. `gotunnel update` and the `go install` launcher verify that signature and the download's checksum before installing, and `gotunnel update` refuses to move to an older version. To check a download yourself, run `minisign -Vm checksums.txt -P RWTnJ6MHF0SvG0vo+EamvWQR10MhqWswNtQ4SqTqG/TEOzSMa6ZcR1Cb`.
- **TCP approvals.** Approving a tunnel now opens only that tunnel. Unapproved attempts on one tunnel no longer cut an IP off from another tunnel that allows it. Approver emails must be plain ASCII.

### Fixed

- **Warning page:** clicking **Continue to site** could show the warning page again instead of opening the site.
- A tunnel could stop accepting new visitors after many visitors closed their connections first. Those slots are now always freed.
- A closed session now always drops its connection, so a client that stops responding after a time limit, kill or ban can't hold its port and capacity slot.
- A client that stops reading can no longer delay its session limit or grow the server's memory use.
- One visitor reading slowly can no longer stall other visitors on the same tunnel. A stream whose visitor stops reading for 5 seconds with its buffers full is closed.
- `--max-bandwidth` and `GOTUNNEL_TCP_MAX_BANDWIDTH` now count every byte, including large frames from modified clients.
- **Self-hosting:** nftables entries for TCP tunnels are now removed when the tunnel ends, so old approvals can't reach a later tunnel on the same port.
- **Self-hosting:** the failed-login count now expires after `--auth-lockout` and is no longer cleared by a successful login from the same IP.

### Changed

- **GoTunnel moved to `gotunnel.site`.** The site is now `gotunnel.site`, docs are at `docs.gotunnel.site`, HTTP tunnels use `*.t.gotunnel.site` and TCP tunnels use `*.tcp.gotunnel.site`. Update with `gotunnel update` to keep using the hosted server.

### Breaking

- **Self-hosting:** `GOTUNNEL_SIGNING_KEY` must now be at least 32 characters. The server refuses to start with a shorter key.
- **Self-hosting:** reload `gotunnel.nft` when you upgrade if you use kernel gating. Bans no longer block IPs a tunnel has approved.
- **Self-hosting:** `--tls-client-ca` now requires `--tls` and a disabled WebSocket listener (`--ws-addr=`), because client certificates can't be checked anywhere else. The server refuses to start otherwise.
- **Self-hosting:** a saved token from before this release is treated as belonging to the hosted server. If you registered with your own server, register again or set `GOTUNNEL_TOKEN`.
- **Stable subdomains and TCP hostnames change once** with this release. Update anything that points at your old address.
- `gotunnel client` with no token and an unencrypted server now stops with instructions instead of starting registration.

---

## [1.2.0] - 2026-10-14

### Changed

- **Self-hosting:** the device token lifetime is now set with `GOTUNNEL_TOKEN_TTL`. The old `GOTUNNEL_BETA_TOKEN_TTL` still works when the new one isn't set.
- **Self-hosting:** the server's registrations table is renamed from `beta_testers` to `devices`. This happens automatically, with no data loss, on the next start or `gotunnel migrate`.

---

## [1.1.0] - 2026-10-14

Install GoTunnel in one line, with npm, or with `go install`.

Versions 1.1.1 to 1.1.3, released the same day, only fixed the npm publishing pipeline and made no changes to GoTunnel itself. If you installed 1.1.1 with npm and have no `gotunnel` command, run `npm install -g gotunnel@latest`.

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
[1.3.0]: https://github.com/bakare-dev/gotunnel-io/releases/tag/v1.3.0
[1.2.0]: https://github.com/bakare-dev/gotunnel-io/releases/tag/v1.2.0
[1.1.0]: https://github.com/bakare-dev/gotunnel-io/releases/tag/v1.1.0
[1.0.0]: https://github.com/bakare-dev/gotunnel-io/releases/tag/v1.0.0
