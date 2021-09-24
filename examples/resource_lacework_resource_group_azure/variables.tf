variable "resource_group_name" {
  type    = string
  default = "Terraform Test Azure Resource Group"
}
variable "description" {
  type    = string
  default = "Terraform Test All Azure Accounts"
}

variable "tenant" {
  type    = string
  default = "a11aa1ab-111a-11ab-a000-11aa1111a11a"
}

variable "subscriptions" {
  type    = list(string)
  default = ["1a1a0b2-abc0-1ab1-1abc-1a000ab0a0a0", "2b000c3-ab10-1a01-1abc-1a000ab0a0a0"]
}
