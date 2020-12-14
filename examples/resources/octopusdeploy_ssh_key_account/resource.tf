resource "octopusdeploy_account" "ssh_key_pair_account" {
  account_type     = "SshKeyPair"
  name             = "SSH Key Pair Account (OK to Delete)"
  private_key_file = "[private_key_file]"
  username         = "[username]"
}