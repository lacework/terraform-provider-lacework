variable "resource_group_name" {
  type = string
  default = "Terraform Test Gcp Resource Group"
}
variable "description" {
  type = string
  default = "Terraform Test All Gcp Accounts"
}

variable "organization" {
  type = string
  default = "Terraform Test All Gcp Accounts"
}

variable "projects" {
  type = list(string)
  default = ["*"]
}
