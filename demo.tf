## Terraform demo project
##
## Usage:
##
## $ tf init
## $ tf apply
## $ tf destroy

variable "n" {
  type = number
  description = "number of items"
  default = 10
}

locals {
  items = [for v in range(var.n) : tostring(v)]
}

resource "time_sleep" "this" {
  for_each = toset(local.items)

  triggers = {
    key = tostring(each.key)
  }
  
  create_duration = "${each.key}s"
  destroy_duration = "${each.key}s"
}

resource "local_file" "this" {
  for_each = toset(local.items)

  content  = "Content ${time_sleep.this[each.key].triggers.key}"
  filename = "./demo-${each.key}.txt"
  file_permission = "0664"
}
