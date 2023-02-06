resource "octopusdeploy_tag_set" "tagset_tag1" {
  name        = "tag1"
  description = "Test tagset"
  sort_order  = 0
}

resource "octopusdeploy_tag" "tag_a" {
  name        = "a"
  color       = "#333333"
  description = "tag a"
  sort_order  = 2
  tag_set_id = octopusdeploy_tag_set.tagset_tag1.id
}

resource "octopusdeploy_tag" "tag_b" {
  name        = "b"
  color       = "#333333"
  description = "tag b"
  sort_order  = 3
  tag_set_id = octopusdeploy_tag_set.tagset_tag1.id
}