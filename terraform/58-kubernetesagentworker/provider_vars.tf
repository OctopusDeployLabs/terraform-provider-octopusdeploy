variable "octopus_server" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "Not used for this test but need to be declared"
}
variable "octopus_apikey" {
  type        = string
  nullable    = false
  sensitive   = true
  description = "Not used for this test but need to be declared"
}
variable "octopus_server_58-kubernetesagentworker" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "The URL of the Octopus server e.g. https://myinstance.octopus.app."
}
variable "octopus_apikey_58-kubernetesagentworker" {
  type        = string
  nullable    = false
  sensitive   = true
  description = "The API key used to access the Octopus server. See https://octopus.com/docs/octopus-rest-api/how-to-create-an-api-key for details on creating an API key."
}
variable "octopus_space_id" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "The space ID to populate"
}
