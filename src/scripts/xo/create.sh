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
eval "$(parseYAML "${YAML}")";

sudo -u postgres psql -c "CREATE USER ${config_db_user}";
sudo -u postgres psql -c "ALTER USER ${config_db_user} WITH PASSWORD '${config_db_pass}'";
sudo -u postgres psql -c "GRANT USAGE, CREATE ON SCHEMA public TO ${config_db_user};";
sudo -u postgres psql -c "CREATE DATABASE ${config_db_name}";
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE ${config_db_name} TO ${config_db_user}";
sudo -u postgres psql -c "ALTER DATABASE ${config_db_name} OWNER TO ${config_db_user};";

echo "Create New DB finished success";