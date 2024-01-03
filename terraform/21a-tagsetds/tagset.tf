data "octopusdeploy_tag_sets" "data_lookup" {
  partial_name = "Test tagset"
  skip         = 0
  take         = 1
}

output "data_lookup" {
  value = data.octopusdeploy_tag_sets.data_lookup.tag_sets[0].id
}

output "tags" {
  value = data.octopusdeploy_tag_sets.data_lookup.tag_sets[0].tags
}