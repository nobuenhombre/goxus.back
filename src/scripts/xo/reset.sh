#!/bin/bash
# Reset goxus_e2e database: drop all tables, run all migrations
# Called before each E2E test run to ensure a clean, reproducible state.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# Load DB config from xo.e2e.yaml (same source as migrate-up-e2e command)
source "$SCRIPT_DIR/yaml.sh"
source "$SCRIPT_DIR/postgresql.sh"

YAML="$1"
eval "$(parseYAML "${YAML}")"

# Write a .pgpass file for passwordless connection
PGPASS_FILE=$(mktemp /tmp/pgpass_e2e_XXXXXX)
chmod 600 "$PGPASS_FILE"
printf '%s:%s:%s:%s:%s\n' "${config_db_host}" "${config_db_port}" "${config_db_name}" "${config_db_user}" "${config_db_pass}" > "$PGPASS_FILE"
trap 'rm -f "$PGPASS_FILE"' EXIT

export PGPASSFILE="$PGPASS_FILE"

echo "=== Resetting ${config_db_name} database ==="

# Step 1: Drop all user tables + the migration version table
echo "  Dropping all tables..."
psql -h "${config_db_host}" -p "${config_db_port}" -U "${config_db_user}" -d "${config_db_name}" -At \
  -c "SELECT 'DROP TABLE IF EXISTS \"' || schemaname || '\".\"' || tablename || '\" CASCADE;'
       FROM pg_tables
       WHERE schemaname = 'public'
         AND tablename != 'spatial_ref_sys';" \
  | psql -h "${config_db_host}" -p "${config_db_port}" -U "${config_db_user}" -d "${config_db_name}"
echo "  Dropped."

# Step 2: Run all migrations (up)
echo "  Running all migrations..."
bash "$SCRIPT_DIR/migrate-up.sh" "$1"
echo "  Done."

echo "=== ${config_db_name} database is fresh ==="
