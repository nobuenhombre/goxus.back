#!/bin/bash
#====================================================#
# Библиотека функций                                 #
#====================================================#

# Функция парсит yaml файл
# и создает локальные переменные
# содержащие значения из yaml файла
#
# пример кода
# # Get YAML file name
# #-------------------
# YAML="$1"
# # Get Variables from YAML file
# #-----------------------------
# eval "$(parseYAML "${YAML}")"
# ...
# ${config_db_host}
#--------------------------------------------
function parseYAML {
   local prefix=${2:-}
   local s='[[:space:]]*' w='[a-zA-Z0-9_]*' fs=$(echo @|tr @ '\034')
   sed -ne "s|^\($s\):|\1|" \
        -e "s|^\($s\)\($w\)$s:$s[\"']\(.*\)[\"']$s\$|\1$fs\2$fs\3|p" \
        -e "s|^\($s\)\($w\)$s:$s\(.*\)$s\$|\1$fs\2$fs\3|p"  $1 |
   awk -F$fs '{
      indent = length($1)/2;
      vname[indent] = $2;
      for (i in vname) {if (i > indent) {delete vname[i]}}
      if (length($3) > 0) {
         vn=""; for (i=0; i<indent; i++) {vn=(vn)(vname[i])("_")}
         printf("%s%s%s=\"%s\"\n", "'$prefix'",vn, $2, $3);
      }
   }'
}