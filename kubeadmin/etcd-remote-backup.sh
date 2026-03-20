#!/usr/bin/env bash
set -Eeuo pipefail

# ============================================================
# Remote etcd backup script
#
# Supports:
#   1. Kubeconfig mode
#   2. Direct etcd cert mode
#   3. CLI args, env vars, and defaults
#
# Example:
#   ./etcd-remote-backup.sh \
#     --endpoint https://192.168.2.32:2379 \
#     --kubeconfig ~/.kube/config
#
# Or:
#   ./etcd-remote-backup.sh \
#     --endpoint https://192.168.2.32:2379 \
#     --cacert /home/ubuntu/ca.crt \
#     --cert /home/ubuntu/healthcheck-client.crt \
#     --key /home/ubuntu/healthcheck-client.key
# ============================================================

# -----------------------------
# DEFAULTS / ENV OVERRIDES
# -----------------------------
ENDPOINT="${ENDPOINT:-}"
KUBECONFIG_PATH="${KUBECONFIG_PATH:-$HOME/.kube/config}"

CACERT="${CACERT:-}"
CERT="${CERT:-}"
KEY="${KEY:-}"

BACKUP_DIR="${BACKUP_DIR:-$HOME/etcd-backups}"
RETENTION_DAYS="${RETENTION_DAYS:-7}"

DIAL_TIMEOUT="${DIAL_TIMEOUT:-5s}"
COMMAND_TIMEOUT="${COMMAND_TIMEOUT:-180s}"

MODE="${MODE:-auto}"   # auto | kubeconfig | certs

TMP_DIR=""
TIMESTAMP="$(date +%F_%H-%M-%S)"
HOSTNAME_SHORT="$(hostname -s)"
SNAPSHOT_FILE=""

# -----------------------------
# HELP
# -----------------------------
usage() {
  cat <<'EOF'
Usage:
  etcd-remote-backup.sh [options]

Options:
  --endpoint URL              etcd endpoint, example: https://192.168.2.32:2379
  --kubeconfig PATH           kubeconfig path (default: ~/.kube/config)
  --cacert PATH               etcd CA certificate path
  --cert PATH                 etcd client certificate path
  --key PATH                  etcd client private key path
  --backup-dir PATH           backup directory
  --retention-days N          delete backups older than N days
  --dial-timeout VALUE        etcdctl dial timeout (default: 5s)
  --command-timeout VALUE     etcdctl command timeout (default: 180s)
  --mode VALUE                auto | kubeconfig | certs
  -h, --help                  show this help

Priority:
  CLI args > environment variables > defaults

Environment variables:
  ENDPOINT
  KUBECONFIG_PATH
  CACERT
  CERT
  KEY
  BACKUP_DIR
  RETENTION_DAYS
  DIAL_TIMEOUT
  COMMAND_TIMEOUT
  MODE

Examples:

  1) Kubeconfig mode:
     ./etcd-remote-backup.sh \
       --endpoint https://192.168.2.32:2379 \
       --kubeconfig ~/.kube/config \
       --mode kubeconfig

  2) Direct cert mode:
     ./etcd-remote-backup.sh \
       --endpoint https://192.168.2.32:2379 \
       --cacert /home/ubuntu/ca.crt \
       --cert /home/ubuntu/healthcheck-client.crt \
       --key /home/ubuntu/healthcheck-client.key \
       --mode certs

  3) Using env vars:
     ENDPOINT=https://192.168.2.32:2379 \
     KUBECONFIG_PATH=$HOME/.kube/config \
     MODE=kubeconfig \
     ./etcd-remote-backup.sh
EOF
}

# -----------------------------
# LOGGING
# -----------------------------
log() {
  echo "[$(date '+%F %T')] $*"
}

fail() {
  echo "[$(date '+%F %T')] ERROR: $*" >&2
  exit 1
}

cleanup() {
  if [[ -n "${TMP_DIR:-}" && -d "${TMP_DIR:-}" ]]; then
    rm -rf "$TMP_DIR"
  fi
}
trap cleanup EXIT

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || fail "Missing command: $1"
}

# -----------------------------
# ARG PARSING
# -----------------------------
parse_args() {
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --endpoint)
        ENDPOINT="$2"
        shift 2
        ;;
      --kubeconfig)
        KUBECONFIG_PATH="$2"
        shift 2
        ;;
      --cacert)
        CACERT="$2"
        shift 2
        ;;
      --cert)
        CERT="$2"
        shift 2
        ;;
      --key)
        KEY="$2"
        shift 2
        ;;
      --backup-dir)
        BACKUP_DIR="$2"
        shift 2
        ;;
      --retention-days)
        RETENTION_DAYS="$2"
        shift 2
        ;;
      --dial-timeout)
        DIAL_TIMEOUT="$2"
        shift 2
        ;;
      --command-timeout)
        COMMAND_TIMEOUT="$2"
        shift 2
        ;;
      --mode)
        MODE="$2"
        shift 2
        ;;
      -h|--help)
        usage
        exit 0
        ;;
      *)
        fail "Unknown argument: $1"
        ;;
    esac
  done
}

# -----------------------------
# KUBECONFIG EXTRACTION
# -----------------------------
extract_certs_from_kubeconfig() {
  [[ -f "$KUBECONFIG_PATH" ]] || fail "Kubeconfig not found: $KUBECONFIG_PATH"

  log "Extracting certs from kubeconfig: $KUBECONFIG_PATH"
  export KUBECONFIG="$KUBECONFIG_PATH"

  TMP_DIR="$(mktemp -d)"

  local ca_data cert_data key_data
  local ca_file cert_file key_file

  ca_data="$(kubectl config view --raw -o jsonpath='{.clusters[0].cluster.certificate-authority-data}' 2>/dev/null || true)"
  cert_data="$(kubectl config view --raw -o jsonpath='{.users[0].user.client-certificate-data}' 2>/dev/null || true)"
  key_data="$(kubectl config view --raw -o jsonpath='{.users[0].user.client-key-data}' 2>/dev/null || true)"

  ca_file="$(kubectl config view --raw -o jsonpath='{.clusters[0].cluster.certificate-authority}' 2>/dev/null || true)"
  cert_file="$(kubectl config view --raw -o jsonpath='{.users[0].user.client-certificate}' 2>/dev/null || true)"
  key_file="$(kubectl config view --raw -o jsonpath='{.users[0].user.client-key}' 2>/dev/null || true)"

  if [[ -n "$ca_data" ]]; then
    echo "$ca_data" | base64 -d > "$TMP_DIR/ca.crt"
    CACERT="$TMP_DIR/ca.crt"
  elif [[ -n "$ca_file" ]]; then
    CACERT="$ca_file"
  else
    fail "No CA cert found in kubeconfig"
  fi

  if [[ -n "$cert_data" ]]; then
    echo "$cert_data" | base64 -d > "$TMP_DIR/client.crt"
    CERT="$TMP_DIR/client.crt"
  elif [[ -n "$cert_file" ]]; then
    CERT="$cert_file"
  else
    fail "No client cert found in kubeconfig"
  fi

  if [[ -n "$key_data" ]]; then
    echo "$key_data" | base64 -d > "$TMP_DIR/client.key"
    KEY="$TMP_DIR/client.key"
  elif [[ -n "$key_file" ]]; then
    KEY="$key_file"
  else
    fail "No client key found in kubeconfig"
  fi
}

# -----------------------------
# MODE RESOLUTION
# -----------------------------
resolve_mode() {
  case "$MODE" in
    certs)
      log "Mode selected: certs"
      ;;
    kubeconfig)
      log "Mode selected: kubeconfig"
      extract_certs_from_kubeconfig
      ;;
    auto)
      if [[ -n "$CACERT" && -n "$CERT" && -n "$KEY" ]]; then
        log "Mode auto -> using direct certs"
      else
        log "Mode auto -> using kubeconfig"
        extract_certs_from_kubeconfig
      fi
      ;;
    *)
      fail "Invalid mode: $MODE (allowed: auto, kubeconfig, certs)"
      ;;
  esac
}

validate_inputs() {
  [[ -n "$ENDPOINT" ]] || fail "ENDPOINT is required. Example: --endpoint https://192.168.2.32:2379"

  [[ -f "$CACERT" ]] || fail "CA cert not found: $CACERT"
  [[ -f "$CERT" ]]   || fail "Client cert not found: $CERT"
  [[ -f "$KEY" ]]    || fail "Client key not found: $KEY"

  mkdir -p "$BACKUP_DIR"
  SNAPSHOT_FILE="${BACKUP_DIR}/etcd-snapshot-${HOSTNAME_SHORT}-${TIMESTAMP}.db"
}

# -----------------------------
# ETCD OPERATIONS
# -----------------------------
etcd_health_check() {
  log "Checking etcd endpoint health..."
  ETCDCTL_API=3 etcdctl \
    --endpoints="$ENDPOINT" \
    --cacert="$CACERT" \
    --cert="$CERT" \
    --key="$KEY" \
    --dial-timeout="$DIAL_TIMEOUT" \
    --command-timeout="10s" \
    endpoint health -w table
}

etcd_endpoint_status() {
  log "Fetching etcd endpoint status..."
  ETCDCTL_API=3 etcdctl \
    --endpoints="$ENDPOINT" \
    --cacert="$CACERT" \
    --cert="$CERT" \
    --key="$KEY" \
    --dial-timeout="$DIAL_TIMEOUT" \
    --command-timeout="15s" \
    endpoint status -w table
}

take_snapshot() {
  log "Starting remote etcd backup..."
  log "Endpoint      : $ENDPOINT"
  log "Backup file   : $SNAPSHOT_FILE"
  log "Dial timeout  : $DIAL_TIMEOUT"
  log "Cmd timeout   : $COMMAND_TIMEOUT"

  ETCDCTL_API=3 etcdctl \
    --endpoints="$ENDPOINT" \
    --cacert="$CACERT" \
    --cert="$CERT" \
    --key="$KEY" \
    --dial-timeout="$DIAL_TIMEOUT" \
    --command-timeout="$COMMAND_TIMEOUT" \
    snapshot save "$SNAPSHOT_FILE"

  log "Backup successful: $SNAPSHOT_FILE"
}

validate_snapshot() {
  log "Validating snapshot..."
  if command -v etcdutl >/dev/null 2>&1; then
    etcdutl snapshot status "$SNAPSHOT_FILE" -w table
  else
    ETCDCTL_API=3 etcdctl snapshot status "$SNAPSHOT_FILE" -w table
  fi
}

cleanup_old_backups() {
  log "Cleaning backups older than ${RETENTION_DAYS} days..."
  find "$BACKUP_DIR" -type f -name "etcd-snapshot-*.db" -mtime +"$RETENTION_DAYS" -delete
}

# -----------------------------
# MAIN
# -----------------------------
main() {
  parse_args "$@"

  require_cmd etcdctl
  require_cmd base64
  require_cmd find

  # kubectl is needed only for kubeconfig mode or auto mode without direct certs
  if [[ "$MODE" == "kubeconfig" ]] || [[ "$MODE" == "auto" && ( -z "$CACERT" || -z "$CERT" || -z "$KEY" ) ]]; then
    require_cmd kubectl
  fi

  resolve_mode
  validate_inputs

  etcd_health_check
  etcd_endpoint_status
  take_snapshot
  validate_snapshot
  cleanup_old_backups

  log "Done."
}

main "$@"