resource "octopusdeploy_certificate" "example" {
  certificate_data = "a-base-64-encoded-string-representing-the-certificate-data"
  name             = "Development Certificate"
  password         = "some-random-value"
}
