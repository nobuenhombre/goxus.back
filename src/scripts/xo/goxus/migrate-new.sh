#!/bin/bash
echo Enter Migration Name [example: create_users_table]?
read migrationName
migrate create -ext sql -dir migrations -seq $migrationName