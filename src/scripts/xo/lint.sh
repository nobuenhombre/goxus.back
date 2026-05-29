#!/bin/bash
source /home/SCRIPTS/__lib__/common.sh
source /home/SCRIPTS/__lib__/yaml.sh
source /home/SCRIPTS/__lib__/go.sh

# Get YAML file name
YAML="$1"

# Get Variables from YAML file
eval "$(parseYAML "${YAML}")"

goFormatCode "${config_codegen_path}"
goLintCode "${config_codegen_path}"

echo "XO finished success"