variable "resource_group_name" {
  type = string
  default = "Terraform Test Machine Resource Group"
}
variable "description" {
  type = string
  default = "Terraform Test All Machine Tags"
}

variable "machine_key" {
  type = string
  default = "test-key"
}

variable "machine_value" {
  type = string
  default = "test-value"
}