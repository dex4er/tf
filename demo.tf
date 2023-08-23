## Terraform demo project
##
## Usage:
##
## $ tf init
## $ tf apply
## $ tf destroy

variable "n" {
  type        = number
  description = "number of items"
  default     = 10
}

locals {
  items = [for v in range(var.n) : "${v}s"]
}

## tflint-ignore: terraform_required_providers
resource "time_sleep" "this" {
  for_each = toset(local.items)

  triggers = {
    key = each.key
  }

  create_duration  = each.key
  destroy_duration = each.key
}

## tflint-ignore: terraform_required_providers
resource "local_file" "this" {
  for_each = toset(local.items)

  content         = "Content ${time_sleep.this[each.key].triggers.key}"
  filename        = "./demo-${each.key}.txt"
  file_permission = "0664"
}

## tflint-ignore: terraform_required_providers,terraform_unused_declarations
data "local_file" "this" {
  for_each = toset(local.items)

  filename = local_file.this[each.key].filename
}

output "filenames" {
  value = [for i in local.items : data.local_file.this[i].filename]
}
