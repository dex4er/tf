#!/bin/sh
set -e
exec terraform init -upgrade \
| grep -v -P 'Finding .* versions matching|Initializing (modules|the backend|provider plugins)...|Upgrading modules...|Using previously-installed|Reusing previous version of|from the shared cache directory|in modules/' \
| sed '/Terraform has been successfully initialized/,$d' \
| uniq \
| sed '1{/^$/d}'
