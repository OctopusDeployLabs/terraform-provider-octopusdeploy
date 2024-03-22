resource "octopusdeploy_machine_policy" "machinepolicy_testing" {
  name                                               = "Testing"
  description                                        = "test machine policy"
  connection_connect_timeout                         = 60000000000
  connection_retry_count_limit                       = 5
  connection_retry_sleep_interval                    = 1000000000
  connection_retry_time_limit                        = 300000000000

  machine_cleanup_policy {
    delete_machines_behavior         = "DeleteUnavailableMachines"
    delete_machines_elapsed_timespan = 1200000000000
  }

  machine_connectivity_policy {
    machine_connectivity_behavior = "ExpectedToBeOnline"
  }

  machine_health_check_policy {

    bash_health_check_policy {
      run_type    = "Inline"
      script_body = ""
    }

    powershell_health_check_policy {
      run_type    = "Inline"
      script_body = "$freeDiskSpaceThreshold = 5GB\r\n\r\nTry {\r\n\tGet-WmiObject win32_LogicalDisk -ErrorAction Stop  | ? { ($_.DriveType -eq 3) -and ($_.FreeSpace -ne $null)} |  % { CheckDriveCapacity @{Name =$_.DeviceId; FreeSpace=$_.FreeSpace} }\r\n} Catch [System.Runtime.InteropServices.COMException] {\r\n\tGet-WmiObject win32_Volume | ? { ($_.DriveType -eq 3) -and ($_.FreeSpace -ne $null) -and ($_.DriveLetter -ne $null)} | % { CheckDriveCapacity @{Name =$_.DriveLetter; FreeSpace=$_.FreeSpace} }\r\n\tGet-WmiObject Win32_MappedLogicalDisk | ? { ($_.FreeSpace -ne $null) -and ($_.DeviceId -ne $null)} | % { CheckDriveCapacity @{Name =$_.DeviceId; FreeSpace=$_.FreeSpace} }\t\r\n}"
    }

    health_check_cron_timezone = "UTC"
    health_check_interval      = 600000000000
    health_check_type          = "RunScript"
  }

  machine_update_policy {
    calamari_update_behavior = "UpdateOnDeployment"
    tentacle_update_behavior = "NeverUpdate"
  }
}