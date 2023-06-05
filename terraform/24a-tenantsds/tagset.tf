data "octopusdeploy_tag_sets" "tagsets" {
  partial_name = "tag1"
  skip = 0
  take = 1
}

output "tagsets" {
  value = data.octopusdeploy_tag_sets.tagsets.tag_sets[0].id
}