resource "octopusdeploy_user" "example" {
  display_name  = "Bob Smith"
  email_address = "bob.smith@example.com"
  is_active     = true
  is_service    = false
  password      = "###########" # get from secure environment/store
  username      = "[username]"

  identity {
    provider = "Octopus ID"
    claim {
      name                 = "email"
      is_identifying_claim = true
      value                = "bob.smith@example.com"
    }
    claim {
      name                 = "dn"
      is_identifying_claim = false
      value                = "Bob Smith"
    }
  }
}
