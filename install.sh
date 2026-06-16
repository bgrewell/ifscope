#!/bin/sh
# ifscope installer / uninstaller.
#
# Install the latest release:
#   curl -fsSL https://raw.githubusercontent.com/bgrewell/ifscope/main/install.sh | sh
#
# Install a specific version, or to a chosen directory:
#   curl -fsSL .../install.sh | sh -s -- --version v0.4.1
#   curl -fsSL .../install.sh | sh -s -- --dir "$HOME/.local/bin"
#
# Uninstall:
#   curl -fsSL .../install.sh | sh -s -- --uninstall
#
# Environment overrides: IFSCOPE_VERSION, IFSCOPE_INSTALL_DIR.
set -eu

REPO="bgrewell/ifscope"
BINARY="ifscope"
VERSION="${IFSCOPE_VERSION:-}"
INSTALL_DIR="${IFSCOPE_INSTALL_DIR:-}"
ACTION="install"

log()  { printf '%s\n' "$*" >&2; }
err()  { printf 'error: %s\n' "$*" >&2; exit 1; }

while [ $# -gt 0 ]; do
	case "$1" in
		--uninstall) ACTION="uninstall" ;;
		--version) VERSION="${2:?--version needs a value}"; shift ;;
		--version=*) VERSION="${1#*=}" ;;
		--dir) INSTALL_DIR="${2:?--dir needs a value}"; shift ;;
		--dir=*) INSTALL_DIR="${1#*=}" ;;
		-h|--help)
			grep '^#' "$0" | sed 's/^# \{0,1\}//'
			exit 0 ;;
		*) err "unknown argument: $1" ;;
	esac
	shift
done

# Pick an install directory: an explicit choice, else /usr/local/bin when
# writable (directly or via sudo), else ~/.local/bin.
SUDO=""
choose_install_dir() {
	[ -n "$INSTALL_DIR" ] && return
	if [ -w /usr/local/bin ]; then
		INSTALL_DIR="/usr/local/bin"
	elif command -v sudo >/dev/null 2>&1; then
		INSTALL_DIR="/usr/local/bin"
		SUDO="sudo"
	else
		INSTALL_DIR="$HOME/.local/bin"
	fi
}

# Run a command with sudo only when the target dir is not user-writable.
as_root() {
	if [ -n "$SUDO" ] || { [ -e "$INSTALL_DIR" ] && [ ! -w "$INSTALL_DIR" ]; }; then
		sudo "$@"
	else
		"$@"
	fi
}

uninstall() {
	target="$INSTALL_DIR/$BINARY"
	if [ ! -f "$target" ]; then
		err "$BINARY not found at $target (use --dir to point at its location)"
	fi
	log "Removing $target"
	as_root rm -f "$target"
	log "ifscope uninstalled."
}

detect_arch() {
	case "$(uname -m)" in
		x86_64|amd64) echo "amd64" ;;
		aarch64|arm64) echo "arm64" ;;
		*) err "unsupported architecture: $(uname -m) (released builds: amd64, arm64)" ;;
	esac
}

# Resolve the latest release tag from the GitHub API without requiring jq.
latest_version() {
	curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" |
		sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p' | head -n 1
}

fetch() { curl -fsSL "$1" -o "$2"; }

install() {
	[ "$(uname -s)" = "Linux" ] || err "ifscope supports Linux only (got $(uname -s))"
	command -v curl >/dev/null 2>&1 || err "curl is required"

	arch="$(detect_arch)"
	if [ -z "$VERSION" ]; then
		log "Resolving latest release..."
		VERSION="$(latest_version)"
		[ -n "$VERSION" ] || err "could not determine latest version; pass --version"
	fi
	ver_no_v="${VERSION#v}"
	archive="${BINARY}_${ver_no_v}_linux_${arch}.tar.gz"
	base="https://github.com/$REPO/releases/download/$VERSION"

	tmp="$(mktemp -d)"
	trap 'rm -rf "$tmp"' EXIT

	log "Downloading $archive ($VERSION)..."
	fetch "$base/$archive" "$tmp/$archive" || err "download failed: $base/$archive"

	# Verify the checksum when sha256sum and the checksums file are available.
	if command -v sha256sum >/dev/null 2>&1 && fetch "$base/checksums.txt" "$tmp/checksums.txt" 2>/dev/null; then
		log "Verifying checksum..."
		( cd "$tmp" && grep " ${archive}\$" checksums.txt | sha256sum -c - ) >/dev/null 2>&1 ||
			err "checksum verification failed for $archive"
	else
		log "warning: skipping checksum verification (sha256sum or checksums.txt unavailable)"
	fi

	tar -xzf "$tmp/$archive" -C "$tmp" "$BINARY" || err "failed to extract $BINARY from archive"
	chmod +x "$tmp/$BINARY"

	choose_install_dir
	[ -d "$INSTALL_DIR" ] || as_root mkdir -p "$INSTALL_DIR"
	as_root mv "$tmp/$BINARY" "$INSTALL_DIR/$BINARY"

	log "Installed $BINARY to $INSTALL_DIR/$BINARY"
	case ":$PATH:" in
		*":$INSTALL_DIR:"*) ;;
		*) log "note: $INSTALL_DIR is not on your PATH; add it to use 'ifscope' directly." ;;
	esac
	"$INSTALL_DIR/$BINARY" --version 2>/dev/null || true
}

case "$ACTION" in
	install)   install ;;
	uninstall) choose_install_dir; uninstall ;;
esac
