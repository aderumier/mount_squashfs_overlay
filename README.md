# squashoverlay

squashoverlay is a drop-in replacement for the `mount.exe` used by
[EmulatorLauncher](https://github.com/RetroBat-Official/emulatorlauncher).
It mounts a read-only squashfs archive as a Windows drive letter,
with a persistent writable overlay directory layered on top using
Docker/OCI-style whiteout conventions (`.wh.<name>` / `.wh..wh..opq`).

## Architecture

- **Lower layer (read-only):** squashfs archive, opened via two libraries:
  - `KarpelesLab/squashfs` — FUSE reads (Open/Stat/ReadDir/Readlink), per-block decompression with ~128 KB RAM per open file.
  - `CalebQ42/squashfs` — copy-on-write materialisation, parallel WriteTo with sync.Pool for efficient bulk decompression.
- **Upper layer (writable):** host directory; receives CoW copies on first write and stores whiteout markers for deleted entries.
- **Case-insensitive path resolution:** squashfs names are matched case-insensitively via a lazily-built, immutably-cached dirIndex.

On Windows the VFS is served through [WinFsp](https://github.com/winfsp/winfsp) (go-winfsp, pure-Go, no CGO).
On Linux a cgofuse/libfuse bridge is used for testing.

## Usage

```
mount.exe [-debug] -drive <X:> [-extractionpath <dir>] [-overlay <dir>] <squashfs-file>
```

| Flag | Description |
|------|-------------|
| `-drive <X:>` | Drive letter to mount at (required) |
| `-overlay <dir>` | Persistent writable overlay directory (omit for read-only mount) |
| `-extractionpath <dir>` | Accepted for compatibility with EmulatorLauncher; ignored |
| `-debug` | Verbose output |

The process runs until killed; killing it unmounts the drive.

### Example

```
mount.exe -drive Z: -extractionpath "C:\Temp\work" -overlay "C:\saves\game1" "C:\roms\game.squashfs"
```

This mounts `game.squashfs` at `Z:`, with any writes or deletions persisted into
`C:\saves\game1\upper\` so they survive remount.

## Requirements

- [WinFsp](https://github.com/winfsp/winfsp/releases) >= 1.10

## Build

```
make
```

Cross-compiles from Linux to Windows with `CGO_ENABLED=0 GOOS=windows GOARCH=amd64`.
Build tags `xz` and `zstd` are enabled automatically (see [Makefile](Makefile)).
