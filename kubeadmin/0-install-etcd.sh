#!/usr/bin/env bash
set -Eeuo pipefail

# -----------------------------
# CONFIG
# -----------------------------

ETCD_VERSION="${ETCD_VERSION:-v3.5.13}"
INSTALL_DIR="/usr/local/bin"
TMP_DIR="/tmp/etcd-download"

ARCH="$(uname -m)"

# Map architecture
case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64 | arm64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

PKG_NAME="etcd-${ETCD_VERSION}-linux-${ARCH}"
DOWNLOAD_URL="https://github.com/etcd-io/etcd/releases/download/${ETCD_VERSION}/${PKG_NAME}.tar.gz"

# -----------------------------
# LOGGING
# -----------------------------

log() {
  echo "[$(date '+%F %T')] $*"
}

fail() {
  echo "ERROR: $*" >&2
  exit 1
}

# -----------------------------
# MAIN
# -----------------------------

main() {
  log "Installing etcdctl version: $ETCD_VERSION"
  log "Architecture: $ARCH"

  mkdir -p "$TMP_DIR"
  cd "$TMP_DIR"

  log "Downloading etcd..."
  wget -q --show-progress "$DOWNLOAD_URL" -O "${PKG_NAME}.tar.gz"

  log "Extracting package..."
  tar -xzf "${PKG_NAME}.tar.gz"

  log "Installing etcdctl to $INSTALL_DIR..."
  sudo mv "${PKG_NAME}/etcdctl" "$INSTALL_DIR/"

  log "Setting permissions..."
  sudo chmod +x "$INSTALL_DIR/etcdctl"

  log "Cleaning up..."
  rm -rf "$TMP_DIR"

  log "Verifying installation..."
  if command -v etcdctl >/dev/null 2>&1; then
    etcdctl version
    log "✅ etcdctl installed successfully"
  else
    fail "etcdctl installation failed"
  fi
}

main "$@"