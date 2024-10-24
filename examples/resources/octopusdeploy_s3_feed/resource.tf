resource "octopusdeploy_s3_feed" "example" {
    name = "AWS S3 Bucket (Ok Delete)"
    use_machine_credentials = false
    access_key = "given_access_key"
    secret_key = "some_secret_key"
}
