# Self-Hosting GoTunnel

GoTunnel ships as a single binary that is **both** the client and the server, so you can run your own tunnel server instead of using the hosted server. Self-hosting needs no database, Redis, or email — that machinery only exists for the hosted server's device enrollment and usage limits. All you need is the binary, a machine with a public IP, and a shared token.

There are two roles:

- **Server** — runs on a machine reachable from the internet (a VPS). `gotunnel server`
- **Client** — runs next to the local service you want to expose. `gotunnel` / `gotunnel client`

## 1. Install the binary on both machines

Install `gotunnel` on the server with any method from the [install guide](INSTALL.md); on Linux the one-liner is `curl -fsSL https://gotunnel.bakaredev.site/install.sh | sh`.

## 2. Start the server

Pick a token your clients will use to authenticate (any long random string):

```bash
# simplest: plaintext control channel, public ports start at 10000
gotunnel server --addr :9000 --start-port 10000 --auth-tokens "your-strong-token"
```

With TLS on the control channel (recommended over the public internet):

```bash
gotunnel server \
  --addr :9000 --tls \
  --tls-cert /path/server-cert.pem \
  --tls-key  /path/server-key.pem \
  --auth-tokens "your-strong-token"
```

- `--auth-tokens` takes one or more comma-separated tokens (or set `GOTUNNEL_AUTH_TOKENS`). There is no default token — the server won't start without one.
- Open the tunnel port (`:9000`) and the public port range (from `--start-port` upward) in your firewall.
- Add `--max-connections`, `--max-auth-failures`, and `--auth-lockout` to taste.

## 3. Connect a client

```bash
export GOTUNNEL_TOKEN="your-strong-token"
gotunnel http 3000 --server your-vps-host:9000
```

You'll get a public `your-vps-host:<assigned-port>` that forwards to `localhost:3000`, with the request inspector on `http://127.0.0.1:4040` like on the hosted server. With TLS on the server, add `--tls --tls-ca /path/ca-cert.pem`, and `--tls-servername` if your certificate name differs from the host you dial.

The token goes in `GOTUNNEL_TOKEN` rather than `--token` so it doesn't show up in the process list. `gotunnel audit` warns if you pass it with `--token`.

## 4. (Optional) Subdomain mode

To hand out `https://<random>.tunnel.example.com` URLs instead of ports, run the server with `--base-domain tunnel.example.com --http-addr 127.0.0.1:8080`, point a wildcard DNS record at the box, and put a TLS-terminating reverse proxy (nginx/Caddy) in front of the HTTP router. This is how the hosted server runs.

## 5. (Optional) TCP tunnels

Let clients share databases and other TCP services with `gotunnel tcp <port>`. Each tunnel gets its own public port from a range you open in your firewall, and is locked to approved IPs.

```bash
GOTUNNEL_TCP_ENABLED=true
GOTUNNEL_TCP_BASE_DOMAIN=tcp.example.com
GOTUNNEL_TCP_PORT_START=20000
GOTUNNEL_TCP_PORT_END=20099
```

- Open the range in your firewall, e.g. `ufw allow 20000:20099/tcp`, and point a **DNS-only** wildcard `*.tcp.example.com` at the server (A record only).
- Without anything else, only `--allowlist` IPs can connect.
- **Approval links:** set `GOTUNNEL_TCP_APPROVAL_URL=https://approve.tcp.example.com`, run the server with `--api-addr`, configure Redis and SMTP for the one-time codes, and proxy `/a/` and `/approve/` on that host to the API address. The approval host must be IPv4-only and DNS-only, because the browser's IP is the one that gets approved.
- **Audit log and bans:** set `GOTUNNEL_DATABASE_URL`.
- **Admin commands:** set `GOTUNNEL_ADMIN_SOCKET=/run/gotunnel/admin.sock`, then `sudo -u <service-user> gotunnel admin sessions`.
- **Kernel gating:** load an nftables table named `gotunnel` with the sets `tcp_allowed`, `tcp_allowed_nets` and `tcp_banned`, grant the service `CAP_NET_ADMIN`, and set `GOTUNNEL_TCP_NFT_ENABLED=true`.

Run `gotunnel help` for every TCP limit (session length, bandwidth, connections per tunnel, rate limits, approved IPs per tunnel).

## Running as a service

Run the server under systemd so it restarts on boot/failure:

```ini
[Unit]
Description=gotunnel
After=network.target

[Service]
Environment=GOTUNNEL_AUTH_TOKENS=your-strong-token
ExecStart=/usr/local/bin/gotunnel server --addr :9000
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

## Updating

```bash
gotunnel update
```

## Hosted vs self-hosted at a glance

| | Hosted server | Self-hosted |
| --- | --- | --- |
| Setup | one command, email verification once | run your own server + token |
| Public URL | `https://<id>.gotunnel.bakaredev.site` | your host/port (or your own subdomain) |
| TLS | handled for you | your certs or a reverse proxy |
| Data path | through the hosted server | entirely your infrastructure |
| Cost | free, with fair-use limits, no guarantees | free, your VPS, your limits |

Both modes speak the same binary protocol with the same stream multiplexing and TLS support — the only difference is who owns the server.
