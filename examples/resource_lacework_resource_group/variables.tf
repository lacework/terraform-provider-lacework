variable "resource_group_name" {
  type    = string
  default = "Terraform Test Aws Resource Group V2"
}
variable "description" {
  type    = string
  default = "Terraform Test RGv2"
}

variable "accounts" {
  type    = list(string)
  default = ["*"]
}
