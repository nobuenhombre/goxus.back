#!/bin/bash
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

source "$SCRIPT_DIR/yaml.sh"
source "$SCRIPT_DIR/postgresql.sh"

YAML="$1"
eval "$(parseYAML "${YAML}")"
CS="$(postgresqlConnectString "postgres" "${config_db_user}" "${config_db_pass}" "${config_db_host}" "${config_db_port}" "${config_db_name}" "${config_db_sslmode}")"

migrate -path "$SCRIPT_DIR/${config_codegen_package}/migrations" -database "${CS}" up
