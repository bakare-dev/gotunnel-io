# Command reference

Every command defaults to the free hosted server, `gotunnel.bakaredev.site`. Pass `--server <host>` to point any of them at a self-hosted server instead. Run `gotunnel help` for the full flag list.

## `gotunnel http <port>`

Expose a local port through the hosted server.

```bash
gotunnel http 3000
gotunnel http 127.0.0.1:8080 --inspect=false
```

Takes a port (`3000` means `localhost:3000`) or a `host:port`. Any client flag can follow it: `--no-reconnect`, `--server`, and so on. The first run on a device asks for your email and registers the device.

### Request inspector

While the tunnel is up, open `http://127.0.0.1:4040` to see every request that came through it. This is handy for testing webhooks. For each request you can:

- see the request and response headers and body, with JSON and form bodies formatted,
- **Replay** it against your local service, without asking the sender to resend it,
- **Copy as cURL** to resend it from a terminal,
- filter the list, or clear it.

The history covers only the current run of `gotunnel`. It is kept in memory on your machine, never written to disk or sent to the GoTunnel server, and disappears when you close the tunnel. It holds the last 200 requests, and up to 64 KB of each request and response. The inspector only listens on localhost. Use `--inspect-addr` to move it, or `--inspect=false` to turn it off.

## `gotunnel tcp <port>`

Expose a database or any other TCP service (Postgres, MySQL, MongoDB, Redis, SSH, RDP). The person connecting uses their normal client and does not install GoTunnel.

```bash
gotunnel tcp 5432
gotunnel tcp 3306 --allowlist 203.0.113.7,10.0.0.0/24
```

You get a host and port, for example `k3x9m2qz7a.tcp.gotunnel.bakaredev.site:20017`, and ready-to-paste commands:

```bash
psql  -h k3x9m2qz7a.tcp.gotunnel.bakaredev.site -p 20017 -U <user> -d <database>
mysql -h k3x9m2qz7a.tcp.gotunnel.bakaredev.site -P 20017 -u <user> -p
```

**TCP tunnels are locked.** A connection is refused unless its IP is approved, in one of two ways:

- **`--allowlist`**: IPs or CIDRs you already know, comma-separated. Ranges wider than `/16` are refused.
- **Approval link**: printed when the tunnel starts. Send it to whoever should connect. They open it on the network they will connect from, verify their email with a one-time code, and their current IP is approved until the tunnel expires. You see each approval live in your terminal.

Limits on the hosted server: one TCP tunnel per account, a fixed session length, up to 5 approved IPs per tunnel, and a cap on connections per tunnel. The port changes every time the tunnel restarts. Your database is still reachable from every approved IP, so keep a strong password and use TLS where your database supports it.

## `gotunnel status`

Shows which server you are pointed at, this device's ID, whether it is registered (email and token expiry), and whether the server is reachable right now.

```text
GoTunnel status
  version        v1.0.0
  server         gotunnel.bakaredev.site (WebSocket over TLS)
  device         6b88fdf5…
  registered     yes
  email          you@gmail.com
  token          valid, expires 2026-11-09 (30 days)
  reachable      yes, 270ms (https://gotunnel.bakaredev.site/enroll/health → 200)
```

Exits with status 1 if the server cannot be reached.

## `gotunnel audit`

Checks how securely this machine connects and prints `PASS`, `WARN`, `FAIL` or `INFO` for each item:

| Check | What it looks at |
| --- | --- |
| transport | Whether the control channel is encrypted (`wss://` or TCP with `--tls`) |
| certificate | The server's TLS certificate: valid, issuer, days until expiry (warns under 14 days) |
| token | Where your token comes from. `--token` on the command line warns, because it lands in shell history; `GOTUNNEL_TOKEN` or device registration pass |
| expiry | Your device token's expiry (warns under 7 days, fails once expired) |
| device | Whether the stored token was issued to this device |
| permissions | That the config folder and token file are readable only by you (skipped on Windows) |

To check a static token without printing it: `gotunnel audit --token "$TOKEN"`. Exits with status 1 if any check fails.

## `gotunnel benchmark`

Measures the network path to the server over fresh connections and reports min, average and p95 for each phase: DNS lookup, TCP connect, TLS handshake, first byte and total.

```bash
gotunnel benchmark          # 5 samples
gotunnel benchmark -n 20    # up to 50
```

It does not open a tunnel or sign in, so it measures how far away the server is, not tunnel throughput. For a raw `host:port` server it times the TCP connect only.

## `gotunnel inspect`

Lists what this build supports (transports, client and server features), its version and platform, and where it keeps its files: the config folder, `token.json` and `device_id`.

This is different from the `--inspect` flag on `gotunnel http`, which starts a local web UI for viewing and replaying requests.

## `gotunnel update`

Replaces the binary with the latest release.

## `gotunnel uninstall`

Removes GoTunnel from this computer. It lists what it will do and asks before doing anything:

1. Frees this device's registration on the server, so your email can be registered again on this or another device.
2. Deletes the config folder, including `token.json` and `device_id`.
3. Deletes the `gotunnel` binary. On macOS and Linux it runs `sudo rm` if the binary is in a folder you can't write to, such as `/usr/local/bin`. On Windows the binary is deleted a moment after the command exits, and the install folder is removed from your user `PATH`.

If the server can't be reached, nothing is removed, so your email isn't left tied to a device you no longer have.

| Flag | What it does |
| --- | --- |
| `--yes` | Skip the confirmation prompt |
| `--force` | Remove GoTunnel even if the server can't free the registration. Your email stays tied to this device |
| `--server` | The server this device registered with (default: the hosted server) |
| `--api-url` | Registration API address, if it isn't `https://<server host>` |

Run it as yourself, not with `sudo`, so it finds your own config folder.

## `gotunnel server`

Runs your own server. See [Self-Hosting](SELF_HOSTING.md).

## `gotunnel admin`

Operator commands, run on the server as the service user. They talk to the running server over a local unix socket (`GOTUNNEL_ADMIN_SOCKET`, owner-only), so they never touch the network.

```bash
sudo -u gotunnel gotunnel admin sessions
sudo -u gotunnel gotunnel admin kill <session-id>
sudo -u gotunnel gotunnel admin ban user@example.com --reason "abuse"
sudo -u gotunnel gotunnel admin unban user@example.com
sudo -u gotunnel gotunnel admin bans
```

- `sessions` lists live HTTP and TCP tunnels, including approved IPs.
- `kill` closes one tunnel.
- `ban` closes every tunnel for the account and its device, and blocks both from connecting and registering again. It needs the database.
