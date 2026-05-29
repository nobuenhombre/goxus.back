#!/bin/bash
source ../yaml.sh
source ../postgresql.sh

eval "$(parseYAML "xo.yaml")"
CS="$(postgresqlConnectString "postgres" "${config_db_user}" "${config_db_pass}" "${config_db_host}" "${config_db_port}" "${config_db_name}" "${config_db_sslmode}")"

migrate -path migrations -database "${CS}" up