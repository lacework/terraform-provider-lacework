variable "integration_name" {
  type    = string
  default = "Github Container Registry Example"
}
variable "username" {
  type      = string
  sensitive = true
}
variable "password" {
  type      = string
  sensitive = true
}
variable "ssl" {
  type = bool
}
variable "non_os_packages" {
  type      = bool
  default   = false
}
