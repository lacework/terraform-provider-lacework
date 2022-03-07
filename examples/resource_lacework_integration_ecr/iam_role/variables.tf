variable "integration_name" {
  type    = string
  default = "Amazon Elastic Container Registry Example"
}
variable "role_arn" {
  type      = string
  sensitive = true
  default   = "arn:aws:iam::1234567890:role/lacework_iam_example_role"
}
variable "external_id" {
  type      = string
  sensitive = true
  default   = "12345"
}
variable "registry_domain" {
  type      = string
  default   = "YourAWSAccount.dkr.ecr.YourRegion.amazonaws.com"
}
variable "non_os_package_support" {
  type      = bool
  default   = true
}
variable "num_images" {
  type      = number
  default   = 10
}

