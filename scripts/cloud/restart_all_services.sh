#!/bin/bash
set -euo pipefail

SERVICES=(
  eigenflux-etcd
  eigenflux-app@profile
  eigenflux-app@item
  eigenflux-app@sort
  eigenflux-app@feed
  eigenflux-app@auth
  eigenflux-app@api
  eigenflux-app@console
  eigenflux-app@pipeline
  eigenflux-app@cron
)

if [[ "${EUID}" -ne 0 ]]; then
  echo "Please run with sudo: sudo ./scripts/cloud/restart_all_services.sh"
  exit 1
fi

echo "Restarting eigenflux services..."

for svc in "${SERVICES[@]}"; do
  echo "==> restarting ${svc}"
  systemctl restart "${svc}"
  systemctl is-active --quiet "${svc}"
  echo "    ${svc} is active"
done

echo ""
echo "All services restarted successfully."
