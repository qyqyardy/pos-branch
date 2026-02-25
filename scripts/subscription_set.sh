#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

if [[ -f ".env" ]]; then
  set -a
  # shellcheck disable=SC1091
  . ".env"
  set +a
fi

PLAN="${1:-}"
PAID_UNTIL_DATE="${2:-}"

if [[ -z "${PLAN}" || -z "${PAID_UNTIL_DATE}" ]]; then
  echo "Usage: $0 <standard|premium> <YYYY-MM-DD>"
  echo "Example: $0 standard 2026-03-25"
  echo ""
  echo "Note: pastikan migration subscription sudah diaplikasikan:"
  echo "  ./scripts/db_apply_migrations.sh"
  exit 1
fi

if [[ "${PLAN}" != "standard" && "${PLAN}" != "premium" ]]; then
  echo "Invalid plan: ${PLAN} (expected: standard|premium)"
  exit 1
fi

DB_CONTAINER="${DB_CONTAINER:-pos_db}"
POSTGRES_USER="${POSTGRES_USER:-pos_user}"
POSTGRES_DB="${POSTGRES_DB:-pos_db}"

if ! docker ps --format '{{.Names}}' | grep -qx "${DB_CONTAINER}"; then
  echo "Container ${DB_CONTAINER} tidak ditemukan / belum jalan."
  echo "Jalankan: docker compose up -d"
  exit 1
fi

PAID_UNTIL_TS="${PAID_UNTIL_DATE} 23:59:59"

echo "Setting subscription: plan=${PLAN}, paid_until=${PAID_UNTIL_TS}"
docker exec -i "${DB_CONTAINER}" psql -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" <<SQL
INSERT INTO store_settings (id, name, tagline, address_line1, address_line2, phone, updated_at, plan, paid_until)
VALUES (1, 'WARKOP', 'Point of Sale', '', '', '', now(), '${PLAN}', '${PAID_UNTIL_TS}')
ON CONFLICT (id) DO UPDATE
SET plan = EXCLUDED.plan,
    paid_until = EXCLUDED.paid_until,
    updated_at = now();
SQL

echo "Done."

