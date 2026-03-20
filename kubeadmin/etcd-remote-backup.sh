#!/usr/bin/env bash
set -Eeuo pipefail

# -----------------------------
# CONFIG (EDIT THESE)
# -----------------------------

# etcd endpoint (use one node only)
ENDPOINT="${ENDPOINT:-https://<CONTROL_PLANE_IP>:2379}"

# Option 1: kubeconfig
KUBECONFIG_PATH="${KUBECONFIG_PATH:-$HOME/admin.conf}"

# Option 2: direct certs (fallback if kubeconfig not used)
CACERT="${CACERT:-}"
CERT="${CERT:-}"
KEY="${KEY:-}"

# Backup settings
BACKUP_DIR="${BACKUP_DIR:-$HOME/etcd-backups}"
RETENTION_DAYS="${RETENTION_DAYS:-7}"

TIMESTAMP="$(date +%F_%H-%M-%S)"
SNAPSHOT_FILE="${BACKUP_DIR}/etcd-snapshot-${TIMESTAMP}.db"

# -----------------------------
# FUNCTIONS
# -----------------------------

log() {
  echo "[$(date '+%F %T')] $*"
}

fail() {
  echo "ERROR: $*" >&2
  exit 1
}

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || fail "Missing command: $1"
}

extract_certs_from_kubeconfig() {
  log "Extracting certs from kubeconfig..."

  TMP_DIR=$(mktemp -d)

  # Extract certs
  CA_DATA=$(kubectl config view --raw -o jsonpath='{.clusters[0].cluster.certificate-authority-data}')
  CERT_DATA=$(kubectl config view --raw -o jsonpath='{.users[0].user.client-certificate-data}')
  KEY_DATA=$(kubectl config view --raw -o jsonpath='{.users[0].user.client-key-data}')

  echo "$CA_DATA"   | base64 -d > "$TMP_DIR/ca.crt"
  echo "$CERT_DATA" | base64 -d > "$TMP_DIR/client.crt"
  echo "$KEY_DATA"  | base64 -d > "$TMP_DIR/client.key"

  CACERT="$TMP_DIR/ca.crt"
  CERT="$TMP_DIR/client.crt"
  KEY="$TMP_DIR/client.key"
}

# -----------------------------
# MAIN
# -----------------------------

main() {
  require_cmd etcdctl
  require_cmd kubectl
  require_cmd base64

  mkdir -p "$BACKUP_DIR"

  export ETCDCTL_API=3

  # If direct certs not provided, extract from kubeconfig
  if [[ -z "$CACERT" || -z "$CERT" || -z "$KEY" ]]; then
    export KUBECONFIG="$KUBECONFIG_PATH"
    extract_certs_from_kubeconfig
  fi

  [[ -f "$CACERT" ]] || fail "CA cert missing"
  [[ -f "$CERT" ]]   || fail "Client cert missing"
  [[ -f "$KEY" ]]    || fail "Client key missing"

  log "Starting remote etcd backup..."
  log "Endpoint: $ENDPOINT"

  etcdctl \
    --endpoints="$ENDPOINT" \
    --cacert="$CACERT" \
    --cert="$CERT" \
    --key="$KEY" \
    snapshot save "$SNAPSHOT_FILE"

  log "Backup successful: $SNAPSHOT_FILE"

  # Validate
  etcdctl snapshot status "$SNAPSHOT_FILE" -w table

  # Cleanup
  find "$BACKUP_DIR" -type f -name "etcd-snapshot-*.db" -mtime +"$RETENTION_DAYS" -delete

  log "Cleanup done (>${RETENTION_DAYS} days)"
}

main "$@"