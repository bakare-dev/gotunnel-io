# GoTunnel Architecture

A conceptual overview of how GoTunnel exposes a local service to the public internet. This document is for people who want to understand how the tool works before trusting it with their traffic. It describes behaviour, not internal source layout.

## High-level overview

GoTunnel is a reverse TCP tunnelling system. A **client** running on your machine keeps a single outbound connection open to a **server** on the public internet. Public visitors reach the server, and the server forwards their traffic down that connection to your local service.

```
+----------------+        +-----------------+        +----------------+
|  Public        |        |  GoTunnel       |        |  GoTunnel      |
|  visitors      | -----> |  server         | <----- |  client        |
|  (browser,     |        |  (public host)  |  one    |  (your machine)|
|   curl, etc.)  |        |                 |  conn   |                |
+----------------+        +-----------------+        +--------+-------+
                                                              |
                                                              v
                                                     Local service
                                                     (localhost:3000, ...)
```

The tunnel itself carries raw bytes. On the hosted server, HTTP tunnels (`gotunnel http`) carry HTTP traffic to a public https URL, and TCP tunnels (`gotunnel tcp`) carry anything else that runs over TCP, such as databases and SSH, locked to approved IPs.

## Three ways to expose a service

| Mode | How a visitor reaches you | When the server uses it |
| --- | --- | --- |
| **Subdomain (HTTP)** | `https://<random>.<base-domain>` routed by the `Host` header | Server started with `--base-domain` (the hosted server uses this) |
| **TCP tunnel** | `<name>.<tcp-base-domain>:<assigned-port>`, routed by port, **locked to approved IPs** | `gotunnel tcp <port>` on a server with `GOTUNNEL_TCP_ENABLED=true` |
| **Port** | `tcp://<server>:<assigned-port>` | Default self-hosted mode (no `--base-domain`) |

Raw TCP carries no hostname, so a TCP tunnel is routed by its own public port from a small range. The name is only a friendly label. Every TCP tunnel starts locked: the server refuses any connection whose source IP has not been approved, before a stream is opened or any byte reaches your machine. An IP is approved in one of two ways:

1. **`--allowlist`** on the owner's command line.
2. **Approval link**: the owner shares a link. The person connecting opens it on the network they will connect from and verifies their email with a one-time code. Their current IP is then approved until the tunnel expires, and the owner sees the approval live.

In subdomain mode the client also connects over **WebSocket (wss) on port 443**, so it works from restrictive networks and reuses a normal TLS certificate. In port mode the client speaks the raw framed protocol directly over TCP (optionally wrapped in TLS).

## Connection lifecycle

1. Server starts and listens for tunnel connections.
2. Client dials the server (with automatic retry and exponential backoff).
3. Client sends a **handshake** declaring its local address and a device id.
4. Server acknowledges the handshake.
5. Client **authenticates** with a token.
6. Server validates the token (constant-time compare against hashed tokens) and, on repeated failures from one IP, temporarily locks that IP out.
7. Server assigns a public endpoint (a port, or a random subdomain) and, **only after the public listener has successfully bound**, replies with `BindOK`.
8. The tunnel is active. Each public connection becomes a multiplexed **stream**.
9. Streams carry data in both directions until closed.
10. Heartbeats keep the session alive; a session that goes silent past the timeout is expired.
11. On Ctrl-C or `q`, the client shuts down cleanly and prints a metrics summary.

## Stream multiplexing

Every public connection maps to one stream, identified by a 32-bit stream id carried in each frame header. Many streams share the single tunnel connection concurrently, and one stream failing does not affect the others. Stream id `0` is reserved for control frames (handshake, auth, heartbeat, and so on).

```
Visitor A ─┐
Visitor B ─┼─► server ──(stream 1,2,3 over one tunnel)──► client ──► local service
Visitor C ─┘
```

## Device registration (hosted server)

The hosted server registers each device once. The first time you connect from a device with no token, the client asks for your Gmail address and the server emails a one-time code; entering it signs the device in. Under the hood:

- The email must be a real `@gmail.com` address, aliases (`+tag`) are rejected, and dots are normalised, so one person maps to one account. Access is limited to **one device per email**.
- On success the server issues an **HMAC-signed, device-bound token** with an expiry. The token is stored on your device and presented on every connect; the tunnel server verifies the signature locally, so normal connections don't hit the accounts database.
- One-time codes live briefly in Redis with a resend cooldown and an attempt limit; enrolled testers are recorded in Postgres.
- `gotunnel uninstall` frees the registration so the email can be used again. The client sends its device ID and current token, and the server deletes the record only if that token matches the one it last issued to the device.

Self-hosters skip all of this and use static tokens (`--auth-tokens` / `GOTUNNEL_AUTH_TOKENS`) instead.

## Security model

GoTunnel is built so that the public-facing surface is the untrusted side.

- **Token authentication** — the server requires one or more tokens (`--auth-tokens` / `GOTUNNEL_AUTH_TOKENS`). Tokens are compared as SHA-256 hashes using a constant-time comparison and are never written to logs. There is no default token; a server with no tokens configured refuses to start.
- **Brute-force protection** — repeated failed authentications from the same IP trigger a temporary lockout (`--max-auth-failures`, `--auth-lockout`).
- **Transport encryption** — TLS 1.2+ is supported on the raw TCP control channel (`--tls`), and the hosted subdomain mode terminates real TLS at the edge. Running without transport encryption is possible for local development but warned against for production.
- **Optional mTLS** — the server can require client certificates signed by a given CA (`--tls-client-ca`), and the client can present a certificate (`--tls-cert` / `--tls-key`) for stronger client identity.
- **Configurable TLS server name** — the client verifies the server certificate against a configurable name (`--tls-servername`), defaulting to the `--server` host.
- **Deadlines and limits** — public-facing reads have a first-byte deadline to blunt slowloris-style abuse, writes have deadlines to avoid hung connections, sessions have a maximum duration and an idle timeout, and the server can cap the number of concurrent tunnels (`--max-connections`).
- **Session isolation** — one client cannot see another client's streams; each session has its own stream namespace, endpoint, and metrics.
- **Frame limits** — a maximum frame size bounds memory per read and rejects oversized payloads.
- **Locked TCP tunnels** — unapproved connections are dropped before a stream opens. Per-IP connection rate limits, a per-tunnel connection cap, one TCP tunnel per account, and a fixed session length bound what an approved IP can do. IPs that repeatedly hit tunnels without approval are temporarily banned.
- **Kernel gating (optional)** — with nftables, the TCP port range is closed in the kernel except for approved (IP, port) pairs, so scanners do not even see the port. Entries expire with the tunnel and are flushed when the server starts.
- **Operator controls** — operators can list and close tunnels and ban accounts through a local, owner-only admin socket. Bans block both connecting and registering. Abuse patterns trigger email alerts.
- **Audit metadata** — for TCP tunnels the server records sessions, approvals and connections (who, when, which IP, how many bytes), never the traffic itself, and deletes them after a retention period (90 days on the hosted server).

## Reliability

- **Heartbeats** run on both ends; a watchdog expires a session that stops responding.
- **Auto-reconnect** on the client re-establishes the tunnel after transient network loss, using exponential backoff with jitter (1s → 30s). It keeps trying until the server is reachable again, unless `--no-reconnect` is set or `--max-retries` caps the attempts.
- **Graceful shutdown** closes public listeners, drains streams, and releases ports on disconnect and on server shutdown.

## Observability

The client shows a live terminal UI with the forwarding endpoint, latency (measured via periodic ping/pong), uptime, TLS/reconnect state, and a colour-coded, timestamped HTTP request log. The server logs session lifecycle events, authentication outcomes (never the token), endpoint assignments, and per-request HTTP summaries. A basic `/healthz` endpoint is available for monitoring.

## Updating

The client can update itself in place:

```bash
gotunnel update
```

It checks the latest published release, downloads the correct build for your platform, and swaps the binary.
