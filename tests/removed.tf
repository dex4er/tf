terraform {
  required_version = ">= 1.7"
}

resource "terraform_data" "this" {
  input = "this"
}
