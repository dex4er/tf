terraform {
  required_version = ">= 1.14"

  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "~> 2.6"
    }
  }
}

action "local_command" "hello" {
  config {
    command   = "echo"
    arguments = ["hello"]
  }
}

resource "terraform_data" "this" {
  input = "this"

  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.local_command.hello]
    }
  }
}
