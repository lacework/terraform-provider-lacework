variable "resource_group_name" {
  type = string
  default = "Terraform Test Container Resource Group"
}
variable "description" {
  type = string
  default = "Terraform Test All Container Tags"
}

variable "ctr_key" {
  type = string
  default = "test-key"
}

variable "ctr_value" {
  type = string
  default = "test-value"
}

variable "ctr_tags" {
  type = list(string)
  default = ["test-tag"]
}