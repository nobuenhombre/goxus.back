#!/bin/bash
source /home/SCRIPTS/__lib__/common.sh
source /home/SCRIPTS/__lib__/yaml.sh
source /home/SCRIPTS/__lib__/ssh.sh
source /home/SCRIPTS/__lib__/postgresql.sh
source /home/SCRIPTS/__lib__/xo.v1.sh
source /home/SCRIPTS/__lib__/go.sh

# Get YAML file name
YAML="$1"

# Get Variables from YAML file
eval "$(parseYAML "${YAML}")"


BACKUP_DATE=$(date --iso-8601)
BACKUP_FILE="${config_db_backups_path}${config_db_name}_${BACKUP_DATE}_dump.sql"
PGPASSWORD=${config_db_pass} psql -h ${config_db_host} -p ${config_db_port} -U ${config_db_user} ${config_db_name} < "${BACKUP_FILE}"


echo "Restore DB from backup finished success"