variable "integration_name" {
  type    = string
  default = "Google Artifact Registry Example"
}
variable "client_id" {
  type      = string
  sensitive = true
}
variable "client_email" {
  type      = string
  sensitive = true
}
variable "private_key_id" {
  type      = string
  sensitive = true
}
variable "private_key" {
  type      = string
  sensitive = true
}
