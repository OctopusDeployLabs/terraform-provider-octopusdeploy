# Basic freeze with no project scopes
resource "octopusdeploy_deployment_freeze" "freeze" {
  name = "Xmas"
  start = "2024-12-25T00:00:00+10:00"
  end = "2024-12-27T00:00:00+08:00"
}

# Freeze with different timezones and single project/environment scope
resource "octopusdeploy_deployment_freeze" "freeze" {
  name = "Xmas"
  start = "2024-12-25T00:00:00+10:00"
  end = "2024-12-27T00:00:00+08:00"
}

# Freeze recurring freeze yearly on Xmas
resource "octopusdeploy_deployment_freeze" "freeze" {
  name = "Xmas"
  start = "2024-12-25T00:00:00+10:00"
  end = "2024-12-27T00:00:00+08:00"
  recurring_schedule = {
    type    = "Annually"
    unit    = 1
    end_type = "Never"
  }
}
