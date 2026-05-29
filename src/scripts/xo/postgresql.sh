#!/bin/bash
#====================================================#
# Библиотека функций                                 #
#====================================================#

# Функция создает строку коннекта к базе PostgreSql
# универсальную
#
# пример кода
# CS="$(postgresqlConnectString "pgsql|postgres" "${config_db_user}" "${config_db_pass}" "${config_db_host}" "${config_db_port}" "${config_db_name}" "${config_db_sslmode}")"
#--------------------------------------------
function postgresqlConnectString {
    local driver=$1
    local user=$2
    local pass=$3
    local host=$4
    local port=$5
    local name=$6
    local sslmode=$7

    local cs="${driver}://${user}:${pass}@${host}:${port}/${name}?sslmode=${sslmode}"

    echo "${cs}"
}

# Функция создает строку коннекта к базе PostgreSql
# для утилиты xo
#
# пример кода
# CS="$(xoPostgresqlConnectString "${config_db_user}" "${config_db_pass}" "${config_db_host}" "${config_db_port}" "${config_db_name}" "${config_db_sslmode}")"
#--------------------------------------------
function xoPostgresqlConnectString {
  local user=$1
  local pass=$2
  local host=$3
  local port=$4
  local name=$5
  local sslmode=$6

  local cs="$(postgresqlConnectString "pgsql" "${user}" "${pass}" "${host}" "${port}" "${name}" "${sslmode}")"

  echo "${cs}"
}

# Функция создает строку коннекта к базе PostgreSql
# для утилиты xouid
#
# пример кода
# CSUID="$(xouidPostgresqlConnectString "${config_db_user}" "${config_db_pass}" "${config_db_host}" "${config_db_port}" "${config_db_name}" "${config_db_sslmode}" "${config_db_pool_max_conns}")"
#--------------------------------------------
function xouidPostgresqlConnectString {
  local user=$1
  local pass=$2
  local host=$3
  local port=$4
  local name=$5
  local sslmode=$6
  local pool_max_conns=$7

  local cs_basic="$(postgresqlConnectString "postgres" "${user}" "${pass}" "${host}" "${port}" "${name}" "${sslmode}")"
  local cs="${cs_basic}&pool_max_conns=${pool_max_conns}"

  echo "${cs}"
}