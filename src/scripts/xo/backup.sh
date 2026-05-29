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

if [[ $config_db_ssh_port ]]; then
  runSshTunnel "${config_db_ssh_port}" "${config_db_ssh_key}" "${config_db_ssh_user}" "${config_db_ssh_host}" "${config_db_host}" "${config_db_port}" "${config_db_ssh_tunnel_db_host}" "${config_db_ssh_tunnel_db_port}"
fi


BACKUP_DATE=$(date --iso-8601)
BACKUP_FILE="${config_db_backups_path}${config_db_name}_${BACKUP_DATE}_dump.sql"
rm -f "${BACKUP_FILE}"
PGPASSWORD=${config_db_pass} pg_dump -h ${config_db_host} -p ${config_db_port} -U ${config_db_user} --dbname=${config_db_name} --file="${BACKUP_FILE}" --create


echo "Backup Production finished success"