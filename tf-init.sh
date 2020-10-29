#!/bin/sh
set -e
exec terraform init \
| grep -v -P '(Initializing (modules|the backend|provider plugins)...|Using previously-installed)' \
| sed '/Terraform has been successfully initialized/,$d' \
| uniq \
| sed '1{/^$/d}'
