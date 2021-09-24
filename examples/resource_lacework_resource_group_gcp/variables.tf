variable "resource_group_name" {
  type    = string
  default = "Terraform Test Gcp Resource Group"
}
variable "description" {
  type    = string
  default = "Terraform Test All Gcp Projects"
}

variable "organization" {
  type    = string
  default = "MyGcpOrg"
}

variable "projects" {
  type    = list(string)
  default = ["*"]
}
