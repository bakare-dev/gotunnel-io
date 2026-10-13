# GoTunnel Protocol Specification

GoTunnel uses a binary, framed, stateful protocol over TCP (optionally TLS) or over a WebSocket carrying binary messages. This document describes protocol **version 1**.

## Transport

- Transport: TCP, optionally wrapped in TLS 1.2+, or a WebSocket (binary frames) for the subdomain/wss path.
- Encoding: binary, big-endian.
- The tunnel payload itself is opaque application data, so any TCP-based protocol is carried unchanged.

## Frame format

Every message is a frame with a fixed 10-byte header followed by a payload.

```
+---------+--------+-----------+-------------+-------------+
| Version | Type   | Stream ID | Payload Len | Payload     |
| 1 byte  | 1 byte | 4 bytes   | 4 bytes     | N bytes     |
+---------+--------+-----------+-------------+-------------+
```

| Field | Size | Description |
| --- | --- | --- |
| Version | 1 byte | Protocol version (`0x01`) |
| Type | 1 byte | Message type (see below) |
| Stream ID | 4 bytes | Stream identifier; `0` for control frames |
| Payload Len | 4 bytes | Payload length in bytes |
| Payload | N bytes | Message payload |

Frames whose declared payload length exceeds the maximum (16 MB) are rejected. Frames with an unknown version are rejected and the connection is closed.

## Message types

| Type | Value | Direction | Description |
| --- | --- | --- | --- |
| `MsgHandshake` | 1 | client → server | Handshake request |
| `MsgHandshakeAck` | 2 | server → client | Handshake acknowledgment |
| `MsgAuth` | 3 | client → server | Authentication token |
| `MsgAuthOK` | 4 | server → client | Authentication accepted |
| `MsgAuthErr` | 5 | server → client | Authentication rejected |
| `MsgBindOK` | 6 | server → client | Public endpoint assigned |
| `MsgStreamOpen` | 7 | server → client | Open a new stream |
| `MsgStreamData` | 8 | both | Stream data |
| `MsgStreamClose` | 9 | both | Close a stream |
| `MsgHeartbeat` | 10 | both | Keepalive |
| `MsgError` | 11 | both | Error / session rejected (payload is a message) |
| `MsgPing` | 12 | client → server | Latency probe (payload: timestamp) |
| `MsgPong` | 13 | server → client | Latency reply (echoes the ping payload) |
| `MsgNotice` | 14 | server → client | Informational notice (JSON), e.g. TCP approval link, IP approved |

Clients ignore message types they do not know, so new server-to-client types are backward compatible.

## Session state machine

The server enforces ordering on inbound frames. Out-of-order frames are rejected and the session is closed.

```
INIT ──MsgHandshake──► HANDSHAKEN ──MsgAuth──► AUTHENTICATED ──► FORWARDING
```

| State | Accepted next | Result |
| --- | --- | --- |
| `INIT` | `MsgHandshake` | → `HANDSHAKEN` (anything else → error) |
| `HANDSHAKEN` | `MsgAuth` | → `AUTHENTICATED` (anything else → error) |
| `AUTHENTICATED` / `FORWARDING` | stream + control frames | stays forwarding |

After authentication the server assigns the public endpoint and sends `MsgBindOK`, moving the connection into active forwarding.

## Flow

### 1. Handshake

Client → server `MsgHandshake`. Payload:

```
+--------+---------------+------------+---------------+-----------+--------------+----------+
| Role   | Capabilities  | ExposeLen  | Expose Addr   | DeviceLen | Device ID    |          |
| 1 byte | 8 bytes       | 2 bytes    | var           | 2 bytes   | var          |          |
+--------+---------------+------------+---------------+-----------+--------------+----------+
```

- **Role**: client (`0x01`) or server (`0x02`).
- **Capabilities**: 64-bit feature bitmask: heartbeat `1<<0`, compression `1<<1`, reconnect `1<<2`, metrics `1<<3`, **TCP `1<<4`** (request a raw TCP tunnel instead of HTTP); remaining bits reserved.
- **Expose Addr**: the client's local service address, e.g. `localhost:3000`.
- **Device ID**: a stable client identifier (used for single-session enforcement).
- **Allowlist** (optional, TCP only): appended after Device ID as a 2-byte length and a comma-separated list of IPs/CIDRs allowed to connect. Omitted entirely when empty, so the wire format of older clients is unchanged. The server rejects invalid entries and ranges wider than `/16` (IPv4) or `/48` (IPv6).

Server → client `MsgHandshakeAck` (no payload).

### 2. Authentication

Client → server `MsgAuth`, payload is the token as raw UTF-8 bytes.

Server → client `MsgAuthOK` (no payload) on success, or `MsgError` with a human-readable message on failure. The token is validated with a constant-time comparison against the server's configured token hashes; it is never logged. Repeated failures from one IP produce a temporary lockout, after which the server replies with `MsgError` describing the retry delay.

### 3. Bind

Server → client `MsgBindOK`, sent only after the public listener has bound. Payload:

```
+-------------+------------------+
| Public Port | Host (optional)  |
| 2 bytes     | var (UTF-8)      |
+-------------+------------------+
```

- **Port mode**: port is the assigned public port, host is empty.
- **Subdomain mode**: port is `0`, host is the assigned `<random>.<base-domain>`.
- **TCP tunnel** (client sent the TCP capability): port is the assigned public TCP port **and** host is `<name>.<tcp-base-domain>`. A client that asked for TCP and gets port `0` or an empty host is talking to a server that cannot lock TCP tunnels, and must disconnect rather than run unlocked.

TCP tunnels are locked: the server refuses every public connection whose source IP is not in the allowlist or approved through the approval link, before any `MsgStreamOpen` is sent.

### 4. Heartbeat and latency

`MsgHeartbeat` (no payload) is sent periodically by both ends; a session that receives no frames within the timeout is expired. The client also sends `MsgPing` carrying an 8-byte timestamp; the server echoes it back as `MsgPong`, letting the client display round-trip latency.

### 5. Stream lifecycle

- **`MsgStreamOpen`** (server → client): a public connection arrived; the stream id in the header identifies it. The client dials the local service.
- **`MsgStreamData`** (both): raw application bytes for that stream, up to 16 MB per frame.
- **`MsgStreamClose`** (both): one side finished; the peer drains and releases the stream.

Each public connection is exactly one stream. Stream ids are assigned sequentially by the server starting at 1.

### 6. Notices

`MsgNotice` (server → client, stream id `0`) carries a small JSON object with a `kind`:

| Kind | Fields | Sent when |
| --- | --- | --- |
| `access` | `approval_url` (empty if the server has no approval page), `expires_at` | Right after `MsgBindOK` for a TCP tunnel |
| `approved` | `email`, `ip` | Someone verified through the approval link and their IP can now connect |

```json
{"kind":"access","approval_url":"https://approve.tcp.example.com/a/7fQx...","expires_at":"2026-10-13T14:30:00Z"}
{"kind":"approved","email":"alice@example.com","ip":"198.51.100.9"}
```

## Limits

| Constraint | Value | Notes |
| --- | --- | --- |
| Max payload size | 16 MB | Per frame |
| Frame header | 10 bytes | Fixed |
| Stream id range | 1 – 2³²−1 | `0` reserved for control |
| Heartbeat interval | 10s | |
| Heartbeat timeout | 30s | Session expires after this |
| Connect timeout | 10s | Initial dial |
| Reconnect backoff | 1s → 30s | Exponential with jitter; retries until reconnected unless `--max-retries` is set |

## Errors

`MsgError` carries a UTF-8 message and is used for authentication failure, lockout, capacity rejection, and session time limits. Invalid frames, unknown versions, oversized payloads, and out-of-order control frames all terminate the session.
