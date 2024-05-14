resource "octopusdeploy_tag_set" "agent_tag_set" {
  name        = "Agent Tag Set"
  description = "Test tagset"
  sort_order  = 0
}

resource "octopusdeploy_tag" "tag_a" {
  name        = "Agent Tag A"
  color       = "#333333"
  description = "tag a"
  sort_order  = 2
  tag_set_id  = octopusdeploy_tag_set.agent_tag_set.id
}

resource "octopusdeploy_tag" "tag_b" {
  name        = "Agent Tag B"
  color       = "#333333"
  description = "tag b"
  sort_order  = 3
  tag_set_id  = octopusdeploy_tag_set.agent_tag_set.id
}

resource "octopusdeploy_tenant" "agent_tenant" {
  name        = "Agent Tenant"
  description = "Test tenant"
  tenant_tags = [octopusdeploy_tag.tag_a.canonical_tag_name, octopusdeploy_tag.tag_b.canonical_tag_name]

  depends_on = [octopusdeploy_tag.tag_a, octopusdeploy_tag.tag_b]
}