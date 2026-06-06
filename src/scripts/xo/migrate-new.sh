#!/bin/bash
DIR="$1"
echo "Enter Migration Name [example: create_users_table]?"
read migrationName
migrate create -ext sql -dir "${DIR}" -seq "${migrationName}"
