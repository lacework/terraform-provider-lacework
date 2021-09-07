variable "resource_group_name" {
  type = string
  default = "Terraform Test Aws Resource Group"
}
variable "description" {
  type = string
  default = "Terraform Test All Aws Accounts"
}

variable "accounts" {
  type = list(string)
  default = ["*"]
}


