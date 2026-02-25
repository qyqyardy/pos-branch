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

DB_CONTAINER="${DB_CONTAINER:-pos_db}"
POSTGRES_USER="${POSTGRES_USER:-pos_user}"
POSTGRES_DB="${POSTGRES_DB:-pos_db}"

if ! docker ps --format '{{.Names}}' | grep -qx "${DB_CONTAINER}"; then
  echo "Container ${DB_CONTAINER} tidak ditemukan / belum jalan."
  echo "Jalankan: docker compose up -d"
  exit 1
fi

shopt -s nullglob
for f in migrations/*.sql; do
  base="$(basename "$f")"
  if [[ "$base" == "001_init.sql" ]]; then
    continue
  fi
  echo "Applying ${base}..."
  docker exec -i "${DB_CONTAINER}" psql -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" < "$f"
done

echo "Migrations applied."
