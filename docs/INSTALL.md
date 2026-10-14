# Download & install

Pick whichever way suits you. They all install the same `gotunnel` binary from the [GitHub releases](https://github.com/bakare-dev/gotunnel-io/releases/latest), and each release lists every binary's SHA-256 in `checksums.txt`.

| Method | Needs | Best for |
| --- | --- | --- |
| [One line, macOS / Linux](#one-line-macos-linux) | Nothing | Most people on macOS or Linux |
| [One line, Windows](#one-line-windows) | PowerShell | Most people on Windows |
| [npm](#npm) | Node.js 22+ | JavaScript developers, or trying it once with `npx` |
| [Go](#go) | Go 1.21+ | Go developers |
| [Manual download](#manual-download) | Nothing | Choosing exactly where the binary goes |

When it's installed, run `gotunnel http 3000` to share a local app. See the [quick start](index.md#quick-start).

## One line: macOS / Linux

```bash
curl -fsSL https://gotunnel.bakaredev.site/install.sh | sh
```

The script:

1. Picks the build for your OS and CPU (Intel or Apple Silicon, x86-64 or ARM64).
2. Downloads it and checks it against the release's `checksums.txt`. If the checksum doesn't match, nothing is installed.
3. Moves it to `/usr/local/bin/gotunnel`. If that folder needs administrator rights, it uses `sudo`, which may ask for your password.

You can [read the script](https://gotunnel.bakaredev.site/install.sh) before running it. It works with `curl` or `wget`.

| Variable | What it does |
| --- | --- |
| `GOTUNNEL_VERSION` | Install a specific release, e.g. `v1.1.0` (default: latest) |
| `GOTUNNEL_INSTALL_DIR` | Install somewhere else (default: `/usr/local/bin`) |

Set them on `sh`, after the pipe. For example, to install without `sudo`:

```bash
curl -fsSL https://gotunnel.bakaredev.site/install.sh | GOTUNNEL_INSTALL_DIR="$HOME/.local/bin" sh
```

## One line: Windows

In PowerShell:

```powershell
irm https://gotunnel.bakaredev.site/install.ps1 | iex
```

The script downloads `gotunnel-windows-amd64.exe`, checks its checksum, installs it as `%LOCALAPPDATA%\Programs\gotunnel\gotunnel.exe` and adds that folder to your user `PATH`. It doesn't need administrator rights. `gotunnel` works in the same window straight away; other windows that were already open need restarting.

You can [read the script](https://gotunnel.bakaredev.site/install.ps1) first. `GOTUNNEL_VERSION` and `GOTUNNEL_INSTALL_DIR` work the same way as on macOS and Linux:

```powershell
$env:GOTUNNEL_VERSION = "v1.1.0"; irm https://gotunnel.bakaredev.site/install.ps1 | iex
```

## npm

Needs Node.js 22 or later.

```bash
npm install -g gotunnel
```

Or run it once without installing:

```bash
npx gotunnel http 3000
```

The package ships the prebuilt binary for your platform in an `@gotunnel/<platform>` package, so nothing is downloaded during install. Supported: macOS (x64, ARM64), Linux (x64, ARM64) and Windows (x64).

If you see `the @gotunnel/… package is missing`, npm skipped optional dependencies. Reinstall without `--omit=optional` or `--no-optional`.

## Go

Needs Go 1.21 or later. Newer Go versions are fetched automatically if needed.

```bash
go install github.com/bakare-dev/gotunnel-io/cmd/gotunnel@latest
```

This installs a small launcher to `$(go env GOPATH)/bin`, usually `~/go/bin`; make sure that folder is on your `PATH`. The first time you run any `gotunnel` command, the launcher downloads the matching release build, checks it against `checksums.txt` and replaces itself with it. After that, `gotunnel` is the normal binary and `gotunnel update` works as usual.

To install a specific release, use its version instead of `latest`, e.g. `@v1.1.0`.

## Manual download

| Platform | Architecture | Download |
| --- | --- | --- |
| macOS | Apple Silicon (M1+) | [gotunnel-darwin-arm64](https://github.com/bakare-dev/gotunnel-io/releases/latest/download/gotunnel-darwin-arm64) |
| macOS | Intel | [gotunnel-darwin-amd64](https://github.com/bakare-dev/gotunnel-io/releases/latest/download/gotunnel-darwin-amd64) |
| Linux | x86-64 | [gotunnel-linux-amd64](https://github.com/bakare-dev/gotunnel-io/releases/latest/download/gotunnel-linux-amd64) |
| Linux | ARM64 | [gotunnel-linux-arm64](https://github.com/bakare-dev/gotunnel-io/releases/latest/download/gotunnel-linux-arm64) |
| Windows | x86-64 | [gotunnel-windows-amd64.exe](https://github.com/bakare-dev/gotunnel-io/releases/latest/download/gotunnel-windows-amd64.exe) |
| All | | [checksums.txt](https://github.com/bakare-dev/gotunnel-io/releases/latest/download/checksums.txt) |

Older versions are on the [releases page](https://github.com/bakare-dev/gotunnel-io/releases).

### macOS and Linux

From the folder you downloaded to (swap in your file name):

```bash
chmod +x gotunnel-darwin-arm64
sudo mv gotunnel-darwin-arm64 /usr/local/bin/gotunnel
gotunnel version
```

On macOS, a binary downloaded with a browser is quarantined, and Gatekeeper may refuse to open it. Clear the flag with:

```bash
xattr -dr com.apple.quarantine /usr/local/bin/gotunnel
```

### Windows

In PowerShell, move the download to a per-user folder and add that folder to your `PATH`:

```powershell
$dir = "$env:LOCALAPPDATA\Programs\gotunnel"
New-Item -ItemType Directory -Force -Path $dir | Out-Null
Move-Item -Force "$HOME\Downloads\gotunnel-windows-amd64.exe" "$dir\gotunnel.exe"
$old = [Environment]::GetEnvironmentVariable("Path","User")
if ($old -notlike "*$dir*") { [Environment]::SetEnvironmentVariable("Path", "$old;$dir", "User") }
$env:Path += ";$dir"
gotunnel version
```

If SmartScreen warns about the file, choose **More info → Run anyway**.

### Check a download yourself

Compare the file's SHA-256 with its line in `checksums.txt`:

```bash
shasum -a 256 gotunnel-darwin-arm64         # macOS
sha256sum gotunnel-linux-amd64              # Linux
```

```powershell
Get-FileHash .\gotunnel-windows-amd64.exe   # Windows
```

## Update

| Installed with | Update with |
| --- | --- |
| One-line script, Go or manual download | `gotunnel update` |
| npm | `npm install -g gotunnel@latest` |

Running the one-line script again also updates. `gotunnel update` on an npm install tells you the npm command instead of changing npm's files.

## Uninstall

```bash
gotunnel uninstall
```

This frees your email on the hosted server so you can register again, deletes your token and device ID, and deletes the binary. On Windows it also removes the install folder from your `PATH`. If you installed with npm, finish with `npm uninstall -g gotunnel`. See [`gotunnel uninstall`](CLI.md#gotunnel-uninstall) for the options.

## Troubleshooting

- **`gotunnel: command not found`.** The install folder isn't on your `PATH`. Open a new terminal, or add the folder the installer printed to your `PATH`. For Go installs, that folder is `$(go env GOPATH)/bin`.
- **An older version runs after installing.** Another `gotunnel` comes first on your `PATH`. `which gotunnel` on macOS/Linux, or `Get-Command gotunnel` on Windows, shows which one.
- **`checksum mismatch`.** The download was corrupted or changed in transit, and nothing was installed. Try again; if it keeps happening, [open an issue](https://github.com/bakare-dev/gotunnel-io/issues/new/choose).
