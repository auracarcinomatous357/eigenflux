#!/usr/bin/env bash

set -euo pipefail

usage() {
  cat <<'EOF'
Usage: ./scripts/cloud/restart.sh <module>

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
  ./scripts/cloud/restart.sh pipeline
  ./scripts/cloud/restart.sh api
  ./scripts/cloud/restart.sh etcd
EOF
}

if [[ $# -ne 1 ]]; then
  usage
  exit 1
fi

module="$1"

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

cmd=(systemctl restart "${service}")
if [[ "${EUID}" -ne 0 ]]; then
  sudo "${cmd[@]}"
  sudo systemctl is-active --quiet "${service}"
else
  "${cmd[@]}"
  systemctl is-active --quiet "${service}"
fi

echo "${service} restarted and active."
