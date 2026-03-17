#!/bin/bash
set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "$0")/../.."; pwd)"

if [[ -f "${PROJECT_ROOT}/.env" ]]; then
  set -a
  # shellcheck disable=SC1091
  source "${PROJECT_ROOT}/.env"
  set +a
fi

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

PORT_CHECKS=(
  "api:${API_PORT:-8080}"
  "console:${CONSOLE_API_PORT:-8090}"
  "profile:${PROFILE_RPC_PORT:-8881}"
  "item:${ITEM_RPC_PORT:-8882}"
  "sort:${SORT_RPC_PORT:-8883}"
  "feed:${FEED_RPC_PORT:-8884}"
  "auth:${AUTH_RPC_PORT:-8886}"
)

status=0

echo "==> systemd services"
for svc in "${SERVICES[@]}"; do
  if systemctl is-active --quiet "${svc}"; then
    echo "[OK]   ${svc}"
  else
    echo "[FAIL] ${svc}"
    status=1
  fi
done

echo ""
echo "==> port listeners"
for entry in "${PORT_CHECKS[@]}"; do
  IFS=':' read -r name port <<< "${entry}"
  if ss -ltn | awk '{print $4}' | grep -Eq "[:.]${port}$"; then
    echo "[OK]   ${name} listening on :${port}"
  else
    echo "[FAIL] ${name} not listening on :${port}"
    status=1
  fi
done

echo ""
echo "==> etcd health"
if docker compose -f "${PROJECT_ROOT}/docker-compose.cloud.yml" ps etcd >/dev/null 2>&1; then
  if docker compose -f "${PROJECT_ROOT}/docker-compose.cloud.yml" exec -T etcd \
    etcdctl --endpoints=http://127.0.0.1:2379 endpoint health >/dev/null 2>&1; then
    echo "[OK]   etcd endpoint healthy"
  else
    echo "[FAIL] etcd endpoint unhealthy"
    status=1
  fi
else
  echo "[FAIL] etcd compose service not found"
  status=1
fi

echo ""
echo "==> quick endpoints"
if curl -fsS "http://127.0.0.1:${API_PORT:-8080}/skill.md" >/dev/null 2>&1; then
  echo "[OK]   api /skill.md"
else
  echo "[FAIL] api /skill.md"
  status=1
fi

if curl -fsS "http://127.0.0.1:${CONSOLE_API_PORT:-8090}/swagger/index.html" >/dev/null 2>&1; then
  echo "[OK]   console /swagger/index.html"
else
  echo "[FAIL] console /swagger/index.html"
  status=1
fi

exit "${status}"
