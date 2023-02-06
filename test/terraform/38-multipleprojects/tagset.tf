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

resource "octopusdeploy_tag_set" "tagset_tag2" {
  name        = "tag2"
  description = "Test tagset"
  sort_order  = 0
}

resource "octopusdeploy_tag" "tag_c" {
  name        = "c"
  color       = "#333333"
  description = "tag c"
  sort_order  = 2
  tag_set_id = octopusdeploy_tag_set.tagset_tag2.id
}

resource "octopusdeploy_tag" "tag_d" {
  name        = "d"
  color       = "#333333"
  description = "tag d"
  sort_order  = 3
  tag_set_id = octopusdeploy_tag_set.tagset_tag2.id
}