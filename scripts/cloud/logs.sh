#!/usr/bin/env bash

set -euo pipefail

usage() {
  cat <<'EOF'
Usage: ./scripts/cloud/logs.sh <module> [journalctl args...]

Modules:
  etcd
  api
  console
  profile
  item
  sort
  feed
  auth
  pipeline
  cron

Examples:
  ./scripts/cloud/logs.sh pipeline
  ./scripts/cloud/logs.sh api -n 100 --no-pager
  ./scripts/cloud/logs.sh etcd --since "10 minutes ago"
EOF
}

if [[ $# -lt 1 ]]; then
  usage
  exit 1
fi

module="$1"
shift

case "${module}" in
  etcd)
    service="eigenflux-etcd"
    ;;
  api|console|profile|item|sort|feed|auth|pipeline|cron)
    service="eigenflux-app@${module}"
    ;;
  *)
    echo "Unknown module: ${module}" >&2
    usage
    exit 1
    ;;
esac

if [[ $# -eq 0 ]]; then
  set -- -f
fi

cmd=(journalctl -u "${service}" "$@")
if [[ "${EUID}" -ne 0 ]]; then
  exec sudo "${cmd[@]}"
fi
exec "${cmd[@]}"
